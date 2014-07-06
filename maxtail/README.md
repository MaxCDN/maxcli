maxtail - alpha
===============

"tail" (sort of) MaxCDN endpoints, return raw json output.

```
Usage: maxtail [arguments...] PATH

Example:

    $ maxtail -a ALIAS -t TOKEN -s SECRET -i 5

Options:

   --config, -c '~/.maxcdn.yml'	yaml file containing all required args
   --alias, -a 			[required] consumer alias
   --token, -t 			[required] consumer token
   --secret, -s 		[required] consumer secret
   --format, -f 'raw'		nginx, raw, json, jsonpp
   --zone, -z 			zone to be tailed (default: all)
   --interval, -i '5'		poll interval in seconds (min: 5)
   --no-follow, -n		print interval and exit
   --quiet, -q			hide 'empty' messages
   --verbose			display verbose http transport information
   --version, -v		print the version
   --help, -h			show help


Formatting Notes:

  "raw"    is basic print of struct.
  "jsonpp" is multiline human readable json output.
  "json"   is single-line parseable json output.
  "nginx"  emulates nginx's default log format...

    Nginx Format
    ------------
    'ClientIp CacheStatus ZoneID [Time] "Uri" '
    'Status Bytes Referer UserAgent OriginTime'

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
```

Installing:
-----------

This can also be installed for system wide use if your `GOBIN` is set via the following:

```bash
# via 'go get' && 'go install'
##

$ go get github.com/MaxCDN/maxcdn-tools/maxtail
$ go install github.com/MaxCDN/maxcdn-tools/maxtail
$ maxtail -h
Usage: maxtail [arguments...] PATH
# ...

# manually
##

git clone https://github.com/MaxCDN/maxcdn-tools
cd maxcdn-tools
make build/maxtail install/maxtail
```
