maxcurl
=======

"curl" (sort of) MaxCDN endpoints, return raw json output.

```
Usage: maxcurl [arguments...] PATH

Example:

    $ maxcurl -a ALIAS -t TOKEN -s SECRET /account.json

Options:

   --config, -c '~/.maxcdn.yml'	yaml file containing all required args
   --alias, -a 			[required] consumer alias
   --token, -t 			[required] consumer token
   --secret, -s 		[required] consumer secret
   --method, -X 'GET'		request method
   --host, -H 			override default API host
   --headers, -i		show headers with body
   --pretty, --pp		pretty print json output
   --verbose			display verbose http transport information
   --version, -v		print the version
   --help, -h			show help



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

Download:
---------

- [linux/386](http://get.maxcdn.com/maxcurl/linux/386/maxcurl) ([md5](http://get.maxcdn.com/maxcurl/linux/386/maxcurl.md5))
- [linux/amd64](http://get.maxcdn.com/maxcurl/linux/amd64/maxcurl) ([md5](http://get.maxcdn.com/maxcurl/linux/amd64/maxcurl.md5))
- [linux/arm](http://get.maxcdn.com/maxcurl/linux/arm/maxcurl) ([md5](http://get.maxcdn.com/maxcurl/linux/arm/maxcurl.md5))
- [darwin/386](http://get.maxcdn.com/maxcurl/darwin/386/maxcurl) ([md5](http://get.maxcdn.com/maxcurl/darwin/amd64/maxcurl.md5))
- [freebsd/386](http://get.maxcdn.com/maxcurl/freebsd/386/maxcurl) ([md5](http://get.maxcdn.com/maxcurl/freebsd/386/maxcurl.md5))
- [freebsd/amd64](http://get.maxcdn.com/maxcurl/freebsd/amd64/maxcurl) ([md5](http://get.maxcdn.com/maxcurl/freebsd/amd64/maxcurl.md5))
- [freebsd/arm](http://get.maxcdn.com/maxcurl/freebsd/arm/maxcurl) ([md5](http://get.maxcdn.com/maxcurl/freebsd/arm/maxcurl.md5))
- [windows/386](http://get.maxcdn.com/maxcurl/windows/386/maxcurl.exe) ([md5](http://get.maxcdn.com/maxcurl/windows/386/maxcurl.exe.md5))
- [windows/amd64](http://get.maxcdn.com/maxcurl/windows/amd64/maxcurl.exe) ([md5](http://get.maxcdn.com/maxcurl/windows/amd64/maxcurl.exe.md5))

Build and Install:
------------------

This can also be installed for system wide use if your `GOBIN` is set via the following:

```bash
# via 'go get' && 'go install'
##

$ go get github.com/MaxCDN/maxcdn-tools/maxcurl
$ go install github.com/MaxCDN/maxcdn-tools/maxcurl
$ maxcurl -h
Usage: maxcurl [arguments...] PATH
# ...

# manually
##

git clone https://github.com/MaxCDN/maxcdn-tools
cd maxcdn-tools
make build/maxcurl install/maxcurl
```
