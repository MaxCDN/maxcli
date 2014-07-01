maxcurl
=======

"curl" (sort of) MaxCDN endpoints, return raw json output.

Installing:
-----------

This can also be installed for system wide use if your `GOBIN` is set via the following:

```bash
# via 'go get' && 'go install'
##

$ go get github.com/jmervine/maxcdn-tools/maxcurl
$ go install github.com/jmervine/maxcdn-tools/maxcurl
$ maxcurl -h
Usage: maxcurl [arguments...] PATH
# ...

# manually
##

git clone https://github.com/jmervine/maxcdn-tools
cd maxcdn-tools
make build/maxcurl install/maxcurl
```
