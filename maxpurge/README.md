maxpurge
========

This provides a simple interface for purging pull zones and their files.

Installing:
-----------

This can also be installed for system wide use if your `GOBIN` is set via the following:

```bash
# via 'go get' && 'go install'
##

$ go get github.com/maxcdn/maxcdn-tools/maxpurge
$ go install github.com/maxcdn/maxcdn-tools/maxpurge
$ maxpurge -h
Usage: maxpurge [arguments...] PATH
# ...

# manually
##

git clone https://github.com/maxcdn/maxcdn-tools
cd maxcdn-tools
make build/maxpurge install/maxpurge
```
