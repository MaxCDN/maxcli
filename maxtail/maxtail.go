package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/MaxCDN/go-maxcdn"
	"github.com/MaxCDN/maxcdn-tools/common"
	"github.com/jmervine/cli"
)

const timeLayout = "2006-01-02T15:04:05"

const (
	name    = "maxtail"
	version = "1.0.0"
)

var config common.Config

func init() {

	// Override cli's default help template
	cli.HelpPrinter = common.HelpPrinter
	cli.VersionPrinter = common.VersionPrinter
	cli.AppHelpTemplate = `Usage: {{.Name}} [arguments...]

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
` + common.AppHelpTemplateFooter

	flags := []cli.Flag{
		cli.StringFlag{"format, f", "raw", "nginx, raw, json, jsonpp"},
		cli.StringFlag{"zone, z", "", "zone to be tailed (default: all)"},
		cli.IntFlag{"interval, i", 60, "poll interval in seconds (min: 5)"},
		cli.BoolFlag{"quiet, q", "hide 'empty' messages"},
		cli.BoolFlag{"no-follow, n", "don't follow, display last 'i' results and exit"},
	}

	actions := func(ctx *cli.Context, cfg *common.Config) {
		if v := ctx.String("format"); v != "" {
			cfg.Format = v
		}

		if v := ctx.String("zone"); v != "" {
			cfg.Zone = v
		}

		interval := ctx.Int("interval")
		if interval < 5 {
			interval = 5
		}

		cfg.Interval = time.Duration(int64(interval) * int64(time.Second))

		cfg.Quiet = ctx.Bool("quiet")
		cfg.NoFollow = ctx.Bool("no-follow")
	}

	config = common.NewCliApp(name, version, flags, actions)
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
		common.Check(err)

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
		if config.NoFollow {
			break
		}
		time.Sleep(config.Interval)
	}
}
