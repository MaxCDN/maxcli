maxreport
=========

Runs reports against [MaxCDN's Reports API](http://docs.maxcdn.com/#reports-api).

Installing:
-----------

This can also be installed for system wide use if your `GOBIN` is set via the following:

```bash
# via 'go get' && 'go install'
##

$ go get github.com/MaxCDN/maxcdn-tools/maxreport
$ go install github.com/MaxCDN/maxcdn-tools/maxreport
$ maxreport -h
Usage: maxreport [arguments...] PATH
# ...

# manually
##

git clone https://github.com/MaxCDN/maxcdn-tools
cd maxcdn-tools
make build/maxreport install/maxreport
```
