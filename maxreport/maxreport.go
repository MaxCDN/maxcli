package main

import (
	"fmt"
	"os"

	"github.com/MaxCDN/go-maxcdn"
	"github.com/MaxCDN/maxcdn-tools/common"
	"github.com/jmervine/cli"
)

const (
	name    = "maxreport"
	version = "1.0.0"
)

var config common.Config

func init() {

	// Override cli's default help template
	cli.HelpPrinter = common.HelpPrinter
	cli.VersionPrinter = common.VersionPrinter
	cli.AppHelpTemplate = `Usage: {{.Name}} [global options] command [command options]

Commands:

    {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
    {{end}}
    For detailed command help, run:

    {{.Name}} command --help

Global Options:

    {{range .Flags}}{{.}}
    {{end}}

` + common.AppHelpTemplateFooter

	app := cli.NewApp()
	app.Name = name
	app.Version = version

	// Setup global flags
	app.Flags = common.DefaultCliFlags

	// Define clobal arguments for inclusion with all commands.
	globals := func(ctx *cli.Context) {
		// Precedence
		// 1. CLI Argument
		// 2. Environment (when applicable)
		// 3. Configuration

		config, _ = common.LoadConfig(ctx.GlobalString("config"))
		common.DefaultCliActions(ctx, &config)
	}

	// Define commands
	app.Commands = []cli.Command{
		{
			Name:        "stats",
			Usage:       "Stats report",
			Description: "Gets the total usage statistics for your account, optionally broken up by {report_type}. If no {report_type} is given the request will return the total usage on your account.",
			Flags: []cli.Flag{
				cli.StringFlag{"from", "", "report start data (YYYY-MM-DD)"},
				cli.StringFlag{"to", "", "report end data (YYYY-MM-DD)"},
				cli.StringFlag{"type, t", "", "report type: hourly, daily, monthly"},
			},
			Action: func(c *cli.Context) {
				globals(c)

				config.Report = "stats"
				config.ReportType = c.String("type")

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
			Usage:       "Popular files report",
			Description: "Gets the most popularly requested files for your account, grouped into daily statistics.",
			Flags: []cli.Flag{
				cli.StringFlag{"from", "", "report start data (YYYY-MM-DD)"},
				cli.StringFlag{"to", "", "report end data (YYYY-MM-DD)"},
				cli.IntFlag{"top, t", 0, "show top N results, zero shows all"},
			},
			Action: func(c *cli.Context) {
				globals(c)

				config.Report = "popular"
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
	common.Check(err)

	stats := data["stats"].(map[string]interface{})
	fmt.Printf("%15s | %15s | %15s | %15s\n", "total hits", "cache hits", "non-cache hits", "size")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("%15v | %15v | %15v | %15v\n", stats["hit"], stats["cache_hit"], stats["noncache_hit"], stats["size"])
	fmt.Println()
}

func statsBreakdown(max *maxcdn.MaxCDN) {
	fmt.Printf("Running %s stats report.\n\n", config.ReportType)

	var data maxcdn.Generic
	endpoint := fmt.Sprintf("/reports/stats.json/%s", config.ReportType)
	_, err := max.Get(&data, endpoint, config.Form)
	common.Check(err)

	fmt.Printf("%25s | %10s | %10s | %10s | %10s\n", "timestamp", "total", "cached", "non-cached", "size")
	fmt.Println(" -------------------------------------------------------------------------------")
	for _, s := range data["stats"].([]interface{}) {
		stats := s.(map[string]interface{})
		fmt.Printf("%25v | %10v | %10v | %10v | %10v\n",
			stats["timestamp"],
			stats["hit"],
			stats["cache_hit"],
			stats["noncache_hit"],
			stats["size"])
	}
	fmt.Println()
}

func popularFiles(max *maxcdn.MaxCDN) {
	fmt.Println("Running popular files report.\n")

	var data maxcdn.Generic
	_, err := max.Get(&data, "/reports/popularfiles.json", config.Form)
	common.Check(err)

	fmt.Printf("%10s | %s\n", "hits", "file")
	fmt.Println("   -----------------")

	for i, f := range data["popularfiles"].([]interface{}) {
		file := f.(map[string]interface{})
		if config.Top != 0 && i == config.Top {
			break
		}
		fmt.Printf("%10v | %v\n", file["hit"], file["uri"])
	}
	fmt.Println()
}
