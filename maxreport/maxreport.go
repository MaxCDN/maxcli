package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/MaxCDN/go-maxcdn"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v1"
)

// Global for configuration.
var config Config

func init() {

	// Override cli's default help template
	cli.AppHelpTemplate = `Usage: {{.Name}} [global options] command [command options]

Commands:
    {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
    {{end}}
    For detailed command help, run:

    {{.Name}} command --help

Global Options:
    {{range .Flags}}{{.}}
    {{end}}
Notes:

    'alias', 'token' and/or 'secret' can be set via exporting them to
    your environment and ALIAS, TOKEN and/or SECRET.

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

`

	// Initialize CLI
	app := cli.NewApp()
	app.Name = "maxreport"
	app.Usage = "Run MaxCDN API Reports"
	app.Version = "1.0.1"
	cli.HelpPrinter = helpPrinter
	cli.VersionPrinter = versionPrinter

	// Setup global flags
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "config, c", Value: "~/.maxcdn.yml", Usage: "yaml file containing all required args"},
		cli.StringFlag{Name: "alias, a", Value: "", Usage: "[required] consumer alias"},
		cli.StringFlag{Name: "token, t", Value: "", Usage: "[required] consumer token"},
		cli.StringFlag{Name: "secret, s", Value: "", Usage: "[required] consumer secret"},
		cli.StringFlag{Name: "host, H", Value: "", Usage: "override default API host"},
		cli.BoolFlag{Name: "verbose", Usage: "display verbose http transport information"},
	}

	// Define clobal arguments for inclusion with all commands.
	globals := func(c *cli.Context) {
		// Precedence
		// 1. CLI Argument
		// 2. Environment (when applicable)
		// 3. Configuration

		config, _ = LoadConfig(c.GlobalString("config"))

		if v := c.GlobalString("alias"); v != "" {
			config.Alias = v
		} else if v := os.Getenv("ALIAS"); v != "" {
			config.Alias = v
		}

		if v := c.GlobalString("token"); v != "" {
			config.Token = v
		} else if v := os.Getenv("TOKEN"); v != "" {
			config.Token = v
		}

		if v := c.GlobalString("secret"); v != "" {
			config.Secret = v
		} else if v := os.Getenv("SECRET"); v != "" {
			config.Secret = v
		}

		if v := config.Validate(); v != "" {
			fmt.Printf("argument error:\n%s\n", v)
			cli.ShowAppHelp(c)
		}

		config.Verbose = c.Bool("verbose")

		if v := c.GlobalString("host"); v != "" {
			config.Host = v
		}
		// handle host override
		if config.Host != "" {
			maxcdn.APIHost = config.Host
		}
	}

	// Define commands
	app.Commands = []cli.Command{
		{
			Name:        "stats",
			Usage:       "stats report",
			Description: "Gets the total usage statistics for your account, optionally broken up by {report_type}. If no {report_type} is given the request will return the total usage on your account.",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "from", Value: "", Usage: "report start data (YYYY-MM-DD)"},
				cli.StringFlag{Name: "to", Value: "", Usage: "report end data (YYYY-MM-DD)"},
				cli.StringFlag{Name: "type, t", Value: "", Usage: "report type: hourly , daily , monthly"},
				cli.BoolFlag{Name: "csv", Usage: "output comma seperated values"},
			},
			Action: func(c *cli.Context) {
				globals(c)

				config.Report = "stats"
				config.ReportType = c.String("type")
				config.CSVOut = c.Bool("csv")

				if f := c.String("from"); f != "" {
					config.Form.Set("from", f)
				}

				if f := c.String("to"); f != "" {
					config.Form.Set("to", f)
				}
			},
		},
		{
			Name:        "popular",
			Usage:       "popular files report",
			Description: "Gets the most popularly requested files for your account, grouped into daily statistics.",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "from", Value: "", Usage: "report start data (YYYY-MM-DD)"},
				cli.StringFlag{Name: "to", Value: "", Usage: "report end data (YYYY-MM-DD)"},
				cli.IntFlag{Name: "top, t", Value: 0, Usage: "show top N results, zero shows all"},
				cli.BoolFlag{Name: "csv", Usage: "output comma seperated values"},
			},
			Action: func(c *cli.Context) {
				globals(c)

				config.Report = "popular"
				config.CSVOut = c.Bool("csv")

				config.Top = c.Int("top")
				if f := c.String("from"); f != "" {
					config.Form.Set("from", f)
				}

				if f := c.String("to"); f != "" {
					config.Form.Set("to", f)
				}
			},
		},
	}

	app.Run(os.Args)
}

func main() {
	max := maxcdn.NewMaxCDN(config.Alias, config.Token, config.Secret)
	max.Verbose = config.Verbose

	switch config.Report {
	case "popular":
		popularFiles(max)
	default:
		stats(max)
	}
}

func stats(max *maxcdn.MaxCDN) {
	if config.ReportType == "" {
		statsSummary(max)
	} else {
		statsBreakdown(max)
	}
}

func statsSummary(max *maxcdn.MaxCDN) {
	fmt.Println("Running summary stats report.\n")

	var data maxcdn.Generic
	_, err := max.Get(&data, "/reports/stats.json", config.Form)
	check(err)

	stats := data["stats"].(map[string]interface{})
	if config.CSVOut {
		statsSummaryCSV(stats)
	} else {
		statsSummaryPrint(stats)
	}
}

func statsSummaryPrint(stats map[string]interface{}) {
	fmt.Printf("%15s | %15s | %15s | %15s\n", "total hits", "cache hits", "non-cache hits", "size")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("%15v | %15v | %15v | %15v\n", stats["hit"], stats["cache_hit"], stats["noncache_hit"], stats["size"])
	fmt.Println()
}

func statsSummaryCSV(stats map[string]interface{}) {
	fmt.Printf("%s,%s,%s,%s\n", "total hits", "cache hits", "non-cache hits", "size")
	fmt.Printf("%v,%v,%v,%v\n", stats["hit"], stats["cache_hit"], stats["noncache_hit"], stats["size"])
}

func statsBreakdown(max *maxcdn.MaxCDN) {
	fmt.Printf("Running %s stats report.\n\n", config.ReportType)

	var data maxcdn.Generic
	endpoint := fmt.Sprintf("/reports/stats.json/%s", config.ReportType)
	_, err := max.Get(&data, endpoint, config.Form)
	check(err)

	stats := data["stats"].([]interface{})
	if config.CSVOut {
		statsBreakdownCSV(stats)
	} else {
		statsBreakdownPrint(stats)
	}
}

func statsBreakdownPrint(stats []interface{}) {
	fmt.Printf("%25s | %10s | %10s | %10s | %10s\n", "timestamp", "total", "cached", "non-cached", "size")
	fmt.Println(" -------------------------------------------------------------------------------")
	for _, s := range stats {
		stats := s.(map[string]interface{})
		fmt.Printf("%25v | %10v | %10v | %10v | %10v\n",
			stats["timestamp"],
			stats["hit"],
			stats["cache_hit"],
			stats["noncache_hit"],
			stats["size"])
	}
}

func statsBreakdownCSV(stats []interface{}) {
	fmt.Printf("%s,%s,%s,%s,%s\n", "timestamp", "total", "cached", "non-cached", "size")
	for _, s := range stats {
		stats := s.(map[string]interface{})
		fmt.Printf("%v,%v,%v,%v,%v\n",
			stats["timestamp"],
			stats["hit"],
			stats["cache_hit"],
			stats["noncache_hit"],
			stats["size"])
	}
}

func popularFiles(max *maxcdn.MaxCDN) {
	fmt.Println("Running popular files report.\n")

	var data maxcdn.Generic
	_, err := max.Get(&data, "/reports/popularfiles.json", config.Form)
	check(err)

	popular := data["popularfiles"].([]interface{})
	if config.CSVOut {
		popularFilesCSV(popular)
	} else {
		popularFilesPrint(popular)
	}
}

func popularFilesPrint(popular []interface{}) {
	fmt.Printf("%10s | %s\n", "hits", "file")
	fmt.Println("   -----------------")

	for i, f := range popular {
		file := f.(map[string]interface{})
		if config.Top != 0 && i == config.Top {
			break
		}
		fmt.Printf("%10v | %v\n", file["hit"], file["uri"])
	}
	fmt.Println()
}

func popularFilesCSV(popular []interface{}) {
	fmt.Printf("%s,%s\n", "hits", "file")

	for i, f := range popular {
		file := f.(map[string]interface{})
		if config.Top != 0 && i == config.Top {
			break
		}
		fmt.Printf("%v,%v\n", file["hit"], file["uri"])
	}
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

/*
 * Config file handlers
 */

type Config struct {
	Host       string `yaml: host,omitempty`
	Alias      string `yaml: alias,omitempty`
	Token      string `yaml: token,omitempty`
	Secret     string `yaml: secret,omitempty`
	Form       url.Values
	Top        int
	Verbose    bool
	Report     string
	ReportType string
	CSVOut     bool
}

func LoadConfig(file string) (c Config, e error) {
	// TODO: this is unix only, look at fixing for windows
	file = strings.Replace(file, "~", os.Getenv("HOME"), 1)

	c = Config{}
	if data, err := ioutil.ReadFile(file); err == nil {
		e = yaml.Unmarshal(data, &c)
	}

	// init empty form, incase we need it
	c.Form = make(url.Values)
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
