maxpurge
========

This provides a simple interface for purging pull zones and their files.

```
Usage: maxpurge [arguments...]
Options:
   --config, -c '~/.maxcdn.yml'			yaml file containing all required args
   --alias, -a 					[required] consumer alias
   --token, -t 					[required] consumer token
   --secret, -s 				[required] consumer secret
   --zone, -z '--zone option --zone option'	[required] zone to be purged
   --file, -f '--file option --file option'	cached file to be purged
   --host, -H 					override default API host
   --verbose					display verbose http transport information
   --version, -v				print the version
   --help, -h					show help


'alias', 'token' 'secret' and/or 'zone' can be set via exporting them
to your environment and ALIAS, TOKEN, SECRET and/or ZONE.

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
    zone: YOUR_ZONE_ID

```

Download:
---------

- [linux/386](http://get.maxcdn.com/maxpurge/linux/386/maxpurge) ([md5](http://get.maxcdn.com/maxpurge/linux/386/maxpurge.md5))
- [linux/amd64](http://get.maxcdn.com/maxpurge/linux/amd64/maxpurge) ([md5](http://get.maxcdn.com/maxpurge/linux/amd64/maxpurge.md5))
- [linux/arm](http://get.maxcdn.com/maxpurge/linux/arm/maxpurge) ([md5](http://get.maxcdn.com/maxpurge/linux/arm/maxpurge.md5))
- [darwin/386](http://get.maxcdn.com/maxpurge/darwin/386/maxpurge) ([md5](http://get.maxcdn.com/maxpurge/darwin/386/maxpurge.md5))
- [darwin/amd64](http://get.maxcdn.com/maxpurge/darwin/amd64/maxpurge) ([md5](http://get.maxcdn.com/maxpurge/darwin/amd64/maxpurge.md5))
- [freebsd/386](http://get.maxcdn.com/maxpurge/freebsd/386/maxpurge) ([md5](http://get.maxcdn.com/maxpurge/freebsd/386/maxpurge.md5))
- [freebsd/amd64](http://get.maxcdn.com/maxpurge/freebsd/amd64/maxpurge) ([md5](http://get.maxcdn.com/maxpurge/freebsd/amd64/maxpurge.md5))
- [freebsd/arm](http://get.maxcdn.com/maxpurge/freebsd/arm/maxpurge) ([md5](http://get.maxcdn.com/maxpurge/freebsd/arm/maxpurge.md5))
- [windows/386](http://get.maxcdn.com/maxpurge/windows/386/maxpurge.exe) ([md5](http://get.maxcdn.com/maxpurge/windows/386/maxpurge.exe.md5))
- [windows/amd64](http://get.maxcdn.com/maxpurge/windows/amd64/maxpurge.exe) ([md5](http://get.maxcdn.com/maxpurge/windows/amd64/maxpurge.exe.md5))


Build and Install:
------------------

This can also be installed for system wide use if your `GOBIN` is set via the following:

```bash
# via 'go get' && 'go install'
##

$ go get github.com/MaxCDN/maxcdn-tools/maxpurge
$ go install github.com/MaxCDN/maxcdn-tools/maxpurge
$ maxpurge -h
Usage: maxpurge [arguments...] PATH
# ...

# manually
##

git clone https://github.com/MaxCDN/maxcdn-tools
cd maxcdn-tools
make build/maxpurge install/maxpurge
```
