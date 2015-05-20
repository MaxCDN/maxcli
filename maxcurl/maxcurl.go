package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/MaxCDN/go-maxcdn"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v1"
)

var config Config

// Full response
type Response struct {
	Code  int
	Data  interface{}
	Error struct {
		Message string
		Type    string
	}
}

func init() {

	// Override cli's default help template
	cli.AppHelpTemplate = `Usage: {{.Name}} [arguments...] PATH

Example:

    $ {{.Name}} -a ALIAS -t TOKEN -s SECRET /account.json

Options:

   {{range .Flags}}{{.}}
   {{end}}


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
	app.Name = "maxcurl"
	app.Version = "1.0.2"

	cli.HelpPrinter = helpPrinter
	cli.VersionPrinter = versionPrinter

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config, c", Value: "~/.maxcdn.yml", Usage: "yaml file containing all required args"},
		cli.StringFlag{Name: "alias, a", Value: "", Usage: "[required] consumer alias"},
		cli.StringFlag{Name: "token, t", Value: "", Usage: "[required] consumer token"},
		cli.StringFlag{Name: "secret, s", Value: "", Usage: "[required] consumer secret"},
		cli.StringFlag{Name: "method, X", Value: "GET", Usage: "request method"},
		cli.StringFlag{Name: "host, H", Value: "", Usage: "override default API host"},
		cli.BoolFlag{Name: "headers, i", Usage: "show headers with body"},
		cli.BoolFlag{Name: "pretty, pp", Usage: "pretty print json output"},
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

		config.Method = c.String("method")
		config.Headers = c.Bool("headers")
		config.Path = c.Args().First()

		if v := config.Validate(); v != "" {
			fmt.Printf("argument error:\n%s\n", v)
			cli.ShowAppHelp(c)
		}

		config.Verbose = c.Bool("verbose")
		if v := c.String("host"); v != "" {
			config.Host = v
		}
		// handle host override
		if config.Host != "" {
			maxcdn.APIHost = config.Host
		}

		if v := c.Bool("pretty"); v {
			config.Pretty = v
		}
	}

	app.Run(os.Args)
}

func main() {
	max := maxcdn.NewMaxCDN(config.Alias, config.Token, config.Secret)
	max.Verbose = config.Verbose

	// seperate path and query
	u, err := url.Parse(config.Path)
	check(err)

	config.Path = u.Path
	form := u.Query()

	// request raw data from maxcdn
	res, err := max.Request(config.Method, config.Path, form)
	defer res.Body.Close()
	check(err)

	body, err := ioutil.ReadAll(res.Body)
	check(err)

	if config.Pretty {
		// format pretty
		var j interface{}
		err = json.Unmarshal(body, &j)
		check(err)

		body, err = json.MarshalIndent(j, "", "  ")
		check(err)
	}

	// print
	if config.Headers {
		fmt.Println(fmtHeaders(res.Header))
	}

	fmt.Printf("%s\n", body)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Replace cli's default help printer with cli's default help printer
// plus an exit at the end.
func helpPrinter(out io.Writer, templ string, data interface{}) {
	w := tabwriter.NewWriter(out, 0, 8, 1, '\t', 0)
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

func fmtHeaders(headers http.Header) (out string) {
	for k, v := range headers {
		out += fmt.Sprintf("%s: %s\n", k, strings.Join(v, ", "))
	}
	return
}

/*
 * Config file handlers
 */

type Config struct {
	Host         string `yaml: host,omitempty`
	Alias        string `yaml: alias,omitempty`
	Token        string `yaml: token,omitempty`
	Secret       string `yaml: secret,omitempty`
	Pretty       bool   `yaml: pretty,omitempty`
	Method, Path string
	Headers      bool
	Verbose      bool
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

	if c.Path == "" {
		out += "- missing path value\n"
	}

	return
}
