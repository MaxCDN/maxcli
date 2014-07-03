package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/MaxCDN/go-maxcdn"
	"github.com/MaxCDN/maxcdn-tools/common"
	"github.com/jmervine/cli"
)

const (
	name    = "maxcurl"
	version = "1.0.1"
)

var config common.Config

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
	cli.HelpPrinter = common.HelpPrinter
	cli.VersionPrinter = common.VersionPrinter
	cli.AppHelpTemplate = `Usage: {{.Name}} [arguments...] PATH

Example:

    $ {{.Name}} -a ALIAS -t TOKEN -s SECRET /account.json

Options:

   {{range .Flags}}{{.}}
   {{end}}

` + common.AppHelpTemplateFooter

	flags := []cli.Flag{
		cli.StringFlag{"method, X", "GET", "request method"},
		cli.BoolFlag{"headers, i", "show headers with body"},
		cli.BoolFlag{"pretty, pp", "pretty print json output"},
	}

	actions := func(ctx *cli.Context, cfg *common.Config) {
		cfg.Validator = customValidator
		cfg.Method = ctx.String("method")
		cfg.Headers = ctx.Bool("headers")
		cfg.Path = ctx.Args().First()

		if v := ctx.Bool("pretty"); v {
			cfg.Pretty = v
		}
	}

	config = common.NewCliApp(name, version, flags, actions)
}

func main() {
	max := maxcdn.NewMaxCDN(config.Alias, config.Token, config.Secret)
	max.Verbose = config.Verbose

	// seperate path and query
	u, err := url.Parse(config.Path)
	common.Check(err)

	config.Path = u.Path
	form := u.Query()

	// request raw data from maxcdn
	res, err := max.Request(config.Method, config.Path, form)
	defer res.Body.Close()
	common.Check(err)

	body, err := ioutil.ReadAll(res.Body)
	common.Check(err)

	if config.Pretty {
		// format pretty
		var j interface{}
		err = json.Unmarshal(body, &j)
		common.Check(err)

		body, err = json.MarshalIndent(j, "", "  ")
		common.Check(err)
	}

	// print
	if config.Headers {
		fmt.Println(fmtHeaders(res.Header))
	}

	fmt.Printf("%s\n", body)
}

func fmtHeaders(headers http.Header) (out string) {
	for k, v := range headers {
		out += fmt.Sprintf("%s: %s\n", k, strings.Join(v, ", "))
	}
	return
}

func customValidator(cfg *common.Config) (out string) {
	if cfg.Alias == "" {
		out += "- missing alias value\n"
	}

	if cfg.Token == "" {
		out += "- missing token value\n"
	}

	if cfg.Secret == "" {
		out += "- missing secret value\n"
	}

	if cfg.Path == "" {
		out += "- missing path value\n"
	}

	return
}
