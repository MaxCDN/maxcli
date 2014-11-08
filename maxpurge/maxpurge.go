package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"text/template"
	"time"

	"github.com/MaxCDN/go-maxcdn"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v1"
)

var start time.Time
var config Config

func init() {
	// Override cli's default help template
	cli.AppHelpTemplate = `Usage: {{.Name}} [arguments...]
Options:
   {{range .Flags}}{{.}}
   {{end}}

'alias', 'token' 'secret' and/or 'zone' can be set via exporting them
to your environment and ALIAS, TOKEN, SECRET and/or ZONE.

Additionally, they can be set in a YAML configuration via the
config option. 'host' can also be set via configuration, but not
environment.

Precedence is argument > environment > configuration.

WARNING:
    Default configuration path works for *nix systems only and
    replies on the 'HOME' environment variable. For Windows, please
    supply a full path.

Sample configuration:

    ---
    alias: YOUR_ALIAS
    token: YOUR_TOKEN
    secret: YOUR_SECRET
    zone: YOUR_ZONE_ID

`

	app := cli.NewApp()

	app.Name = "maxpurge"
	app.Version = "1.0.0"

	cli.HelpPrinter = helpPrinter
	cli.VersionPrinter = versionPrinter

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config, c", Value: "~/.maxcdn.yml", Usage: "yaml file containing all required args"},
		cli.StringFlag{Name: "alias, a", Value: "", Usage: "[required] consumer alias"},
		cli.StringFlag{Name: "token, t", Value: "", Usage: "[required] consumer token"},
		cli.StringFlag{Name: "secret, s", Value: "", Usage: "[required] consumer secret"},
		cli.IntSliceFlag{Name: "zone, z", Value: new(cli.IntSlice), Usage: "[required] zone to be purged"},
		cli.StringSliceFlag{Name: "file, f", Value: new(cli.StringSlice), Usage: "cached file to be purged"},
		cli.StringFlag{Name: "host, H", Value: "", Usage: "override default API host"},
		cli.BoolFlag{Name: "verbose", Usage: "display verbose http transport information"},
	}

	app.Action = func(c *cli.Context) {
		// Precedence
		// 1. CLI Argument
		// 2. Environment (when applicable)
		// 3. Configuration

		config, _ = LoadConfig(c.String("config"))

		if v := c.String("alias"); v != "" {
			config.Alias = v
		} else if v := os.Getenv("ALIAS"); v != "" {
			config.Alias = v
		}

		if v := c.String("token"); v != "" {
			config.Token = v
		} else if v := os.Getenv("TOKEN"); v != "" {
			config.Token = v
		}

		if v := c.String("secret"); v != "" {
			config.Secret = v
		} else if v := os.Getenv("SECRET"); v != "" {
			config.Secret = v
		}

		if v := c.IntSlice("zone"); len(v) != 0 {
			config.Zones = v
		} else if v := os.Getenv("ZONE"); v != "" {
			zones := strings.Split(v, ",")
			for i, z := range zones {
				n, err := strconv.ParseInt(strings.TrimSpace(z), 0, 64)
				check(err)

				config.Zones[i] = int(n)
			}
		}

		config.Files = c.StringSlice("file")
		config.Verbose = c.Bool("verbose")

		if v := config.Validate(); v != "" {
			fmt.Printf("argument error:\n%s\n", v)
			cli.ShowAppHelp(c)
		}

		if v := c.String("host"); v != "" {
			config.Host = v
		}
		// handle host override
		if config.Host != "" {
			maxcdn.APIHost = config.Host
		}
	}

	app.Run(os.Args)

	start = time.Now()
}

func main() {
	max := maxcdn.NewMaxCDN(config.Alias, config.Token, config.Secret)
	max.Verbose = config.Verbose

	var response []*maxcdn.Response
	var err error
	var successful bool

	if len(config.Files) != 0 {
		var resps []*maxcdn.Response
		for _, zone := range config.Zones {
			resps, err = max.PurgeFiles(zone, config.Files)
			response = append(response, resps...)
		}
		successful = len(response) == (len(config.Files) * len(config.Zones))
	} else {
		response, err = max.PurgeZones(config.Zones)
		successful = (len(response) == len(config.Zones))
	}
	check(err)

	if successful {
		fmt.Printf("Purge successful after: %v.\n", time.Since(start))
	} else {
		check(fmt.Errorf("error: one or more of your purges did not succeed"))
	}
}

func check(err error) {
	if err != nil {
		fmt.Printf("%v\n\nPurge failed after %v.\n", err, time.Since(start))
		os.Exit(2)
	}
}

// Replace cli's default help printer with cli's default help printer
// plus an exit at the end.
func helpPrinter(templ string, data interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	t := template.Must(template.New("help").Parse(templ))
	err := t.Execute(w, data)
	if err != nil {
		panic(err)
	}
	w.Flush()
	os.Exit(0)
}

// Replace cli's default version printer with cli's default version printer
// plus an exit at the end.
func versionPrinter(c *cli.Context) {
	fmt.Printf("%v version %v\n", c.App.Name, c.App.Version)
	os.Exit(0)
}

/*
 * Config file handlers
 */

type Config struct {
	Host    string `yaml: host,omitempty`
	Alias   string `yaml: alias,omitempty`
	Token   string `yaml: token,omitempty`
	Secret  string `yaml: secret,omitempty`
	Zones   []int  `yaml: secret,omitempty`
	Files   []string
	Verbose bool
}

func LoadConfig(file string) (c Config, e error) {
	// TODO: this is unix only, look at fixing for windows
	file = strings.Replace(file, "~", os.Getenv("HOME"), 1)

	c = Config{}
	if data, err := ioutil.ReadFile(file); err == nil {
		e = yaml.Unmarshal(data, &c)
	}

	return
}

func (c *Config) Validate() (out string) {
	if c.Alias == "" {
		out += "- missing alias value\n"
	}

	if c.Token == "" {
		out += "- missing token value\n"
	}

	if c.Secret == "" {
		out += "- missing secret value\n"
	}

	if len(c.Zones) == 0 {
		out += "- missing zones value\n"
	}

	return
}
