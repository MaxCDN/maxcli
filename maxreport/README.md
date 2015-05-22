maxreport
=========

Runs reports against [MaxCDN's Reports API](http://docs.maxcdn.com/#reports-api).


##### Global Help

```
$ maxreport --help
Usage: maxreport [global options] command [command options]

Commands:
    stats	stats report
    popular	popular files report
    help, h	Shows a list of commands or help for one command

    For detailed command help, run:

    maxreport command --help

Global Options:
    --config, -c '~/.maxcdn.yml'	yaml file containing all required args
    --alias, -a 			[required] consumer alias
    --token, -t 			[required] consumer token
    --secret, -s 			[required] consumer secret
    --host, -H 				override default API host
    --verbose				display verbose http transport information
    --version, -v			print the version
    --help, -h				show help

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
```

##### Stats Help
```
$ maxreport stats --help
NAME:
   stats - stats report

USAGE:
   command stats [command options] [arguments...]

DESCRIPTION:
   Gets the total usage statistics for your account, optionally broken up by {report_type}. If no {report_type} is given the request will return the total usage on your account.

OPTIONS:
   --from 	report start data (YYYY-MM-DD)
   --to 	report end data (YYYY-MM-DD)
   --type, -t 	report type: hourly, daily, monthly
   --csv	output comma seperated values
```

##### Popular Help
```
$ maxreport popular --help
NAME:
   popular - popular files report

USAGE:
   command popular [command options] [arguments...]

DESCRIPTION:
   Gets the most popularly requested files for your account, grouped into daily statistics.

OPTIONS:
   --from 		report start data (YYYY-MM-DD)
   --to 		report end data (YYYY-MM-DD)
   --top, -t '0'	show top N results, zero shows all
   --csv		output comma seperated values
```

Download:
---------

- [linux-386](http://get.maxcdn.com/maxreport/linux/386/maxreport) ([md5](http://get.maxcdn.com/maxreport/linux/386/maxreport.md5))
- [linux-amd64](http://get.maxcdn.com/maxreport/linux/amd64/maxreport) ([md5](http://get.maxcdn.com/maxreport/linux/amd64/maxreport.md5))
- [linux-arm](http://get.maxcdn.com/maxreport/linux/arm/maxreport) ([md5](http://get.maxcdn.com/maxreport/linux/arm/maxreport.md5))
- [darwin-386](http://get.maxcdn.com/maxreport/darwin/386/maxreport) ([md5](http://get.maxcdn.com/maxreport/darwin/386/maxreport.md5))
- [darwin-amd64](http://get.maxcdn.com/maxreport/darwin/amd64/maxreport) ([md5](http://get.maxcdn.com/maxreport/darwin/amd64/maxreport.md5))
- [freebsd-386](http://get.maxcdn.com/maxreport/freebsd/386/maxreport) ([md5](http://get.maxcdn.com/maxreport/freebsd/386/maxreport.md5))
- [freebsd-amd64](http://get.maxcdn.com/maxreport/freebsd/amd64/maxreport) ([md5](http://get.maxcdn.com/maxreport/freebsd/amd64/maxreport.md5))
- [freebsd-arm](http://get.maxcdn.com/maxreport/freebsd/arm/maxreport) ([md5](http://get.maxcdn.com/maxreport/freebsd/arm/maxreport.md5))
- [windows-386](http://get.maxcdn.com/maxreport/windows/386/maxreport.exe) ([md5](http://get.maxcdn.com/maxreport/windows/386/maxreport.exe.md5))
- [windows-amd64](http://get.maxcdn.com/maxreport/windows/amd64/maxreport.exe) ([md5](http://get.maxcdn.com/maxreport/windows/amd64/maxreport.exe.md5))



Build and Install:
------------------

This can also be installed for system wide use if your `GOBIN` is set via the following:

```bash
# via 'go get' && 'go install'
##

$ go get github.com/MaxCDN/maxcli/maxreport
$ go install github.com/MaxCDN/maxcli/maxreport
$ maxreport -h
Usage: maxreport [arguments...] PATH
# ...

# manually
##

git clone https://github.com/MaxCDN/maxcli
cd maxcli
make build/maxreport install/maxreport
```
