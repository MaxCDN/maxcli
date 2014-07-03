maxtail
=======

"tail" (sort of) MaxCDN endpoints, return raw json output.

### Work in progress.

**This currently isn't working as intended.**

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
