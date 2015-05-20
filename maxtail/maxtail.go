package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"
	"time"

	"github.com/MaxCDN/go-maxcdn"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v1"
)

const timeLayout = "2006-01-02T15:04:05"
const timeDelay = 15 * time.Second
const minInterval = 5

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

- "raw"    is basic print of struct.
- "jsonpp" is multiline human readable json output.
- "json"   is single-line parseable json output.
- "nginx"  emulates nginx's default log format...

    Nginx Format
    ------------
    'ClientIp CacheStatus ZoneID [Time] "Uri" '
    'Status Bytes Referer UserAgent OriginTime'

Filter Notes:

- "zones":
    The specific zones whose requests you want to pull.  Separate
    multiple zone ids by comma.
- "uri":
    Use this filter to view requests made for a specific resource
    (or group of resources). You can do a literal match or regular
    expression in this field (i.e. '/images/header.png' or
    'regex:/images/').
- "status":
    The specific HTTP status code responses you want to pull.
    Separate multiple HTTP status codes by comma (i.e. 200,201,304).
- "ssl":
    Use this filter to distinguish between SSL and non-SSL traffic
    (choose nossl, ssl or both).
- "ua":
    Filter logs by specific user agents. You can do a literal match
    or regular expression in this field (i.e. 'Python MaxCDN API
    Client' or 'regex:Chrome').
- "referer":
    Filter logs by a specific referer. You can do a literal match or
    regular expression in this field (i.e. 'www.maxcdn.com' or
    'regex:maxcdn.com').
- "pop":
    Filter logs by specific POPs (Points Of Presence), use comma
    separation for multiple POPs.
- "qs":
    Filter logs by a specific query string. You can do a literal
    match or regular expression in this field (i.e. 'width=600' or
    'regex:width').

See https://docs.maxcdn.com/#get-raw-logs for full details.

Alpha Notes:

- There is a built in 15 second delay on events.
- Does not currently show all events while in follow mode. However,
  the more you filters applied, the closer it will be to displaying
  all events.

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
	app.Version = "0.0.4-alpha"

	cli.HelpPrinter = helpPrinter
	cli.VersionPrinter = versionPrinter

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config, c", Value: "~/.maxcdn.yml", Usage: "yaml file containing all required args"},
		cli.StringFlag{Name: "alias, a", Value: "", Usage: "[required] consumer alias"},
		cli.StringFlag{Name: "token, t", Value: "", Usage: "[required] consumer token"},
		cli.StringFlag{Name: "secret, s", Value: "", Usage: "[required] consumer secret"},
		cli.StringFlag{Name: "format, f", Value: "raw", Usage: "nginx, raw, json, jsonpp"},
		cli.IntFlag{Name: "interval, i", Value: 5, Usage: "poll interval in seconds (min: 5)"},
		cli.BoolFlag{Name: "no-follow, n", Usage: "print interval and exit"},
		cli.BoolFlag{Name: "quiet, q", Usage: "hide 'empty' messages"},
		cli.BoolFlag{Name: "verbose", Usage: "display verbose http transport information"},

		// Filters
		cli.StringFlag{Name: "zones", Value: "", Usage: "filter: by zone(s)"},
		cli.StringFlag{Name: "uri", Value: "", Usage: "filter: uri"},
		cli.StringFlag{Name: "status", Value: "", Usage: "filter: status"},
		cli.StringFlag{Name: "ssl", Value: "", Usage: "filter: ssl"},
		cli.StringFlag{Name: "ua", Value: "", Usage: "filter: user agent"},
		cli.StringFlag{Name: "referer", Value: "", Usage: "filter: referer"},
		cli.StringFlag{Name: "pop", Value: "", Usage: "filter: pop"},
		cli.StringFlag{Name: "qs", Value: "", Usage: "filter: query string"},
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

		if v := c.String("zones"); v != "" {
			config.ZonesFilter = v
		}

		if v := c.String("uri"); v != "" {
			config.UriFilter = v
		}

		if v := c.String("status"); v != "" {
			config.StatusFilter = v
		}

		if v := c.String("ssl"); v != "" {
			config.SSLFilter = v
		}

		if v := c.String("ua"); v != "" {
			config.UAFilter = v
		}

		if v := c.String("referer"); v != "" {
			config.RefFilter = v
		}

		if v := c.String("pop"); v != "" {
			config.POPFilter = v
		}

		if v := c.String("qs"); v != "" {
			config.QSFilter = v
		}

		interval := c.Int("interval")
		if interval < minInterval {
			interval = minInterval
		}

		config.Interval = time.Duration(int64(interval) * int64(time.Second))

		config.NoFollow = c.Bool("no-follow")
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

	for {
		adj := config.Interval + timeDelay
		now := time.Now().UTC()
		start := now.Add(-adj).Format(timeLayout)
		end := now.Add(-timeDelay).Format(timeLayout)

		form := url.Values{}
		form.Set("start", start)
		form.Set("end", end)
		form.Set("sort", "oldest")

		if config.ZonesFilter != "" {
			form.Set("zones", config.ZonesFilter)
		}

		if config.UriFilter != "" {
			form.Set("uri", config.UriFilter)
		}

		if config.SSLFilter != "" {
			form.Set("ssl", config.SSLFilter)
		}

		if config.StatusFilter != "" {
			form.Set("status", config.StatusFilter)
		}

		if config.UAFilter != "" {
			form.Set("user_agent", config.UAFilter)
		}

		if config.POPFilter != "" {
			form.Set("pop", config.POPFilter)
		}

		if config.QSFilter != "" {
			form.Set("query_string", config.QSFilter)
		}

		if config.RefFilter != "" {
			form.Set("referer", config.RefFilter)
		}

		logs, err := max.GetLogs(form)
		check(err)

		if len(logs.Records) == 0 {
			if !config.Quiet {
				log.Println("empty ... ")
			}
		} else {

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
					/*
					   Nginx Format
					   ------------
					   'ClientIp CacheStatus ZoneID [Time] "Uri" '
					   'Status Bytes Referer UserAgent OriginTime'
					*/
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

		if config.NoFollow {
			break
		}

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
	Interval time.Duration
	Quiet    bool
	Verbose  bool
	NoFollow bool

	// Filters
	ZonesFilter  string
	UriFilter    string
	SSLFilter    string
	StatusFilter string
	UAFilter     string
	POPFilter    string
	QSFilter     string
	RefFilter    string
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
