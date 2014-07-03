package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/MaxCDN/go-maxcdn"
	"github.com/MaxCDN/maxcdn-tools/common"
	"github.com/jmervine/cli"
)

const (
	name    = "maxpurge"
	version = "1.0.1"
)

var start time.Time
var config common.Config

func init() {
	// Override cli's default help template
	cli.HelpPrinter = common.HelpPrinter
	cli.VersionPrinter = common.VersionPrinter
	cli.AppHelpTemplate = `Usage: {{.Name}} [arguments...]
Options:
   {{range .Flags}}{{.}}
   {{end}}
` + common.AppHelpTemplateFooter

	flags := []cli.Flag{
		cli.IntSliceFlag{"zone, z", new(cli.IntSlice), "[required] zone to be purged"},
		cli.StringSliceFlag{"file, f", new(cli.StringSlice), "cached file to be purged"},
	}

	actions := func(ctx *cli.Context, cfg *common.Config) {
		cfg.Validator = customValidator
		if v := ctx.IntSlice("zone"); len(v) != 0 {
			cfg.Zones = v
		} else if v := os.Getenv("ZONE"); v != "" {
			zones := strings.Split(v, ",")
			for i, z := range zones {
				n, err := strconv.ParseInt(strings.TrimSpace(z), 0, 64)
				check(err)

				cfg.Zones[i] = int(n)
			}
		}

		cfg.Files = ctx.StringSlice("file")
	}

	config = common.NewCliApp(name, version, flags, actions)
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

	if len(cfg.Zones) == 0 {
		out += "- missing zone value\n"
	}

	return
}
