package common

import (
	"fmt"
	"github.com/MaxCDN/go-maxcdn"
	"github.com/jmervine/cli"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"
	"time"
)

// usage

var AppHelpTemplateFooter = `Credential Notes:

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

var DefaultCliFlags = []cli.Flag{
	cli.StringFlag{"config, c", "~/.maxcdn.yml", "yaml file containing all required args"},
	cli.StringFlag{"alias, a", "", "[required] consumer alias"},
	cli.StringFlag{"token, t", "", "[required] consumer token"},
	cli.StringFlag{"secret, s", "", "[required] consumer secret"},
	cli.BoolFlag{"verbose", "display verbose http transport information"},
	cli.StringFlag{"host, H", "", "override default API host"},
}

var DefaultCliActions = func(ctx *cli.Context, cfg *Config) {
	if v := ctx.String("alias"); v != "" {
		cfg.Alias = v
	} else if v := os.Getenv("ALIAS"); v != "" {
		cfg.Alias = v
	}

	if v := ctx.String("token"); v != "" {
		cfg.Token = v
	} else if v := os.Getenv("TOKEN"); v != "" {
		cfg.Token = v
	}

	if v := ctx.String("secret"); v != "" {
		cfg.Secret = v
	} else if v := os.Getenv("SECRET"); v != "" {
		cfg.Secret = v
	}

	cfg.Verbose = ctx.Bool("verbose")

	if v := cfg.Validator(cfg); v != "" {
		fmt.Printf("argument error:\n%s\n", v)
		cli.ShowAppHelp(ctx)
	}

	if v := ctx.String("host"); v != "" {
		cfg.Host = v
	}
	// handle host override
	if cfg.Host != "" {
		maxcdn.APIHost = cfg.Host
	}
}

func NewCliApp(name, version string, flags []cli.Flag, actions func(c *cli.Context, cfg *Config)) Config {
	var cfg Config
	app := cli.NewApp()
	app.Name = name
	app.Version = version

	app.Flags = DefaultCliFlags
	app.Flags = append(app.Flags, flags...)

	app.Action = func(ctx *cli.Context) {
		// Precedence
		// 1. CLI Argument
		// 2. Environment (when applicable)
		// 3. Configuration

		cfg, _ = LoadConfig(ctx.String("config"))

		actions(ctx, &cfg)
		DefaultCliActions(ctx, &cfg)
	}

	app.Run(os.Args)
	return cfg
}

type Config struct {
	Host    string `yaml: host,omitempty`
	Alias   string `yaml: alias,omitempty`
	Token   string `yaml: token,omitempty`
	Secret  string `yaml: secret,omitempty`
	Verbose bool

	Validator func(*Config) string

	// maxtail
	Format   string `yaml: format,omitempty`
	Zone     string `yaml: zone,omitempty`
	Interval time.Duration
	Quiet    bool
	NoFollow bool

	// maxpurge
	Zones []int `yaml: secret,omitempty`
	Files []string

	// maxcurl
	Pretty  bool `yaml: pretty,omitempty`
	Method  string
	Headers bool
	Path    string

	// maxreport
	Form       url.Values
	Top        int
	Report     string
	ReportType string
}

func LoadConfig(file string) (c Config, e error) {
	// TODO: this is unix only, look at fixing for windows
	file = strings.Replace(file, "~", os.Getenv("HOME"), 1)

	c = Config{}
	c.Validator = defaultValidator
	if data, err := ioutil.ReadFile(file); err == nil {
		e = yaml.Unmarshal(data, &c)
	}

	return
}

func defaultValidator(cfg *Config) (out string) {
	if cfg.Alias == "" {
		out += "- missing alias value\n"
	}

	if cfg.Token == "" {
		out += "- missing token value\n"
	}

	if cfg.Secret == "" {
		out += "- missing secret value\n"
	}
	return
}

// Replace cli's default help printer with cli's default help printer
// plus an exit at the end.
func HelpPrinter(templ string, data interface{}) {
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
func VersionPrinter(c *cli.Context) {
	fmt.Printf("%v version %v\n", c.App.Name, c.App.Version)
	os.Exit(0)
}

func Check(err error) {
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}
