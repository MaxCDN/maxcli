package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"
	"time"

	"github.com/MaxCDN/go-maxcdn"
	"github.com/jmervine/cli"
	"gopkg.in/yaml.v1"
)

const timeLayout = "2006-01-02T15:04:05"

var config Config

func init() {

	// Override cli's default help template
	cli.AppHelpTemplate = `Usage: {{.Name}} [arguments...] PATH

Example:

    $ {{.Name}} -a ALIAS -t TOKEN -s SECRET -i 5

Options:

   {{range .Flags}}{{.}}
   {{end}}

Formatting Notes:

  "raw"    is basic print of struct.
  "jsonpp" is multiline human readable json output.
  "json"   is single-line parseable json output.
  "nginx"  emulates nginx's default log format...

    Nginx Format
    ------------
    'ClientIp CacheStatus ZoneID [Time] "Uri" '
    'Status Bytes Referer UserAgent OriginTime'

Credential Notes:

'alias', 'token' and/or 'secret' can be set via exporting them to
your environment and ALIAS, TOKEN and/or SECRET.

Additionally, they can be set in a YAML configuration via the
config option. 'pretty' and 'host' can also be set via
configuration, but not environment.

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
    pretty: true

`

	app := cli.NewApp()
	app.Name = "maxtail"
	app.Version = "1.0.0"

	cli.HelpPrinter = helpPrinter
	cli.VersionPrinter = versionPrinter

	app.Flags = []cli.Flag{
		cli.StringFlag{"config, c", "~/.maxcdn.yml", "yaml file containing all required args"},
		cli.StringFlag{"alias, a", "", "[required] consumer alias"},
		cli.StringFlag{"token, t", "", "[required] consumer token"},
		cli.StringFlag{"secret, s", "", "[required] consumer secret"},
		cli.StringFlag{"format, f", "raw", "nginx, raw, json, jsonpp"},
		cli.StringFlag{"zone, z", "", "zone to be tailed (default: all)"},
		cli.IntFlag{"interval, i", 60, "poll interval in seconds (min: 5)"},
		cli.BoolFlag{"quiet, q", "hide 'empty' messages"},
		cli.BoolFlag{"verbose", "display verbose http transport information"},
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

		if v := c.String("format"); v != "" {
			config.Format = v
		}

		if v := c.String("zone"); v != "" {
			config.Zone = v
		}

		interval := c.Int("interval")
		if interval < 5 {
			interval = 5
		}

		config.Interval = time.Duration(int64(interval) * int64(time.Second))

		config.Quiet = c.Bool("quiet")

		config.Verbose = c.Bool("verbose")
		if v := c.String("host"); v != "" {
			config.Host = v
		}

		// handle host override
		if config.Host != "" {
			maxcdn.APIHost = config.Host
		}
	}

	app.Run(os.Args)
}

func main() {
	max := maxcdn.NewMaxCDN(config.Alias, config.Token, config.Secret)
	max.Verbose = config.Verbose

	tailer := func(s, e string) {
		form := url.Values{}
		form.Set("start", s)
		form.Set("end", e)

		if config.Zone != "" {
			form.Set("zone", config.Zone)
		}

		logs, err := max.GetLogs(form)
		check(err)

		if len(logs.Records) == 0 {
			if !config.Quiet {
				log.Println("empty ... ")
			}
			return
		}

		/*
		   Nginx Format
		   ------------
		   'ClientIp CacheStatus ZoneID [Time] "Uri" '
		   'Status Bytes Referer UserAgent OriginTime'
		*/

		for _, line := range logs.Records {
			switch config.Format {
			case "jsonpp":
				s, e := json.MarshalIndent(line, "", "\t")
				if e != nil {
					log.Printf("%+v\n", e)
				} else {
					log.Printf("\n-------------------------\n%s\n\n", s)
				}
			case "json":
				s, e := json.Marshal(line)
				if e != nil {
					log.Printf("%+v\n", e)
				} else {
					fmt.Printf("%s\n", s)
				}
			case "nginx":
				fmt.Printf("%s %s %d [%s] %q %d %d %q %q %.3f\n",
					line.ClientIp,
					line.CacheStatus,
					line.ZoneID,
					line.Time,
					line.Uri,
					line.Status,
					line.Bytes,
					line.Referer,
					line.UserAgent,
					line.OriginTime)
			default:
				log.Printf("%+v\n", line)
			}

		}
	}

	end := time.Now()
	start := end.Add(-config.Interval)

	for {
		end = start
		start = start.Add(-config.Interval)

		tailer(start.Format(timeLayout), end.Format(timeLayout))
		time.Sleep(config.Interval)
	}
}

/*
 * Config file handlers
 */

type Config struct {
	Host     string `yaml: host,omitempty`
	Alias    string `yaml: alias,omitempty`
	Token    string `yaml: token,omitempty`
	Secret   string `yaml: secret,omitempty`
	Format   string `yaml: format,omitempty`
	Zone     string `yaml: zone,omitempty`
	Interval time.Duration
	Quiet    bool
	Verbose  bool
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

	return
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Replace cli's default help printer with cli's default help printer
// plus an exit at the end.
func helpPrinter(templ string, data interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	t := template.Must(template.New("help").Parse(templ))
	err := t.Execute(w, data)
	check(err)

	w.Flush()
	os.Exit(0)
}

// Replace cli's default version printer with cli's default version printer
// plus an exit at the end.
func versionPrinter(c *cli.Context) {
	fmt.Printf("%v version %v\n", c.App.Name, c.App.Version)
	os.Exit(0)
}
