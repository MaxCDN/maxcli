maxcdn tools
============

Collection of CLI tools for interfacing with MaxCDN's REST API

> Built using [github.com/MaxCDN/go-maxcdn](https://github.com/MaxCDN/go-maxcdn).

#### Note

`maxtail` currently isn't working as intended. See [issue #2](https://github.com/MaxCDN/maxcdn-tools/issues/2) for updates.

Configuration
-------------

All tools use a configuration file as it's last means of getting `alias`, `secret` and
`token`. See individal tool `help` for addtional options available in your configuration.

```yaml
---
alias: YOUR_ALIAS
token: YOUR_TOKEN
secret: YOUR_SECRET
```

See [sample.maxcdn.yml](sample.maxcdn.yml) for a more complete example.


Installing:
-----------

This can also be installed for system wide use if your `GOBIN` is set via the following:

```bash
# via 'go get' && 'go install'
##

go get github.com/MaxCDN/maxcdn-tools/{{tool}}
go install github.com/MaxCDN/maxcdn-tools/{{tool}}

# manually
##

git clone https://github.com/MaxCDN/maxcdn-tools
cd maxcdn-tools

# build and install all tools
make build install

# or a single tool
make build/{{tool}} install/{{tool}}

# windows users
cd {{tool}}
go build
move {{tool}} c:\path\to\GOBIN
```

See individal tool README for additional instructions.

Prebuilt Binaries:
------------------

A set of binaries for all tools have been prebuilt using golang's cross compiler on `Linux 3.8.0-36-generic #52~precise1-Ubuntu SMP x86_64`.

**Building All Binaries:**

```bash
make setup # only once, this takes a while
make build/all

# or
make build/all/{{tool}}
```

Here's what's available for each tool:

- maxreport
    - [linux-386](http://get.maxcdn.com/maxreport/linux/386/maxpurge) ([md5](http://get.maxcdn.com/maxreport/linux/386/maxpurge.md5))
    - [linux-amd64](http://get.maxcdn.com/maxreport/linux/amd64/maxpurge) ([md5](http://get.maxcdn.com/maxreport/linux/amd64/maxpurge.md5))
    - [linux-arm](http://get.maxcdn.com/maxreport/linux/arm/maxpurge) ([md5](http://get.maxcdn.com/maxreport/linux/arm/maxpurge.md5))
    - [darwin-386](http://get.maxcdn.com/maxreport/darwin/386/maxpurge) ([md5](http://get.maxcdn.com/maxreport/darwin/386/maxpurge.md5))
    - [darwin-amd64](http://get.maxcdn.com/maxreport/darwin/amd64/maxpurge) ([md5](http://get.maxcdn.com/maxreport/darwin/amd64/maxpurge.md5))
    - [freebsd-386](http://get.maxcdn.com/maxreport/freebsd/386/maxpurge) ([md5](http://get.maxcdn.com/maxreport/freebsd/386/maxpurge.md5))
    - [freebsd-amd64](http://get.maxcdn.com/maxreport/freebsd/amd64/maxpurge) ([md5](http://get.maxcdn.com/maxreport/freebsd/amd64/maxpurge.md5))
    - [freebsd-arm](http://get.maxcdn.com/maxreport/freebsd/arm/maxpurge) ([md5](http://get.maxcdn.com/maxreport/freebsd/arm/maxpurge.md5))
    - [windows-386](http://get.maxcdn.com/maxreport/windows/386/maxpurge.exe) ([md5](http://get.maxcdn.com/maxreport/windows/386/maxpurge.exe.md5))
    - [windows-amd64](http://get.maxcdn.com/maxreport/windows/amd64/maxpurge.exe) ([md5](http://get.maxcdn.com/maxreport/windows/amd64/maxpurge.exe.md5))
- maxpurge
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
- maxcurl
    - [linux/386](http://get.maxcdn.com/maxcurl/linux/386/maxcurl) ([md5](http://get.maxcdn.com/maxcurl/linux/386/maxcurl.md5))
    - [linux/amd64](http://get.maxcdn.com/maxcurl/linux/amd64/maxcurl) ([md5](http://get.maxcdn.com/maxcurl/linux/amd64/maxcurl.md5))
    - [linux/arm](http://get.maxcdn.com/maxcurl/linux/arm/maxcurl) ([md5](http://get.maxcdn.com/maxcurl/linux/arm/maxcurl.md5))
    - [darwin/386](http://get.maxcdn.com/maxcurl/darwin/386/maxcurl) ([md5](http://get.maxcdn.com/maxcurl/darwin/amd64/maxcurl.md5))
    - [freebsd/386](http://get.maxcdn.com/maxcurl/freebsd/386/maxcurl) ([md5](http://get.maxcdn.com/maxcurl/freebsd/386/maxcurl.md5))
    - [freebsd/amd64](http://get.maxcdn.com/maxcurl/freebsd/amd64/maxcurl) ([md5](http://get.maxcdn.com/maxcurl/freebsd/amd64/maxcurl.md5))
    - [freebsd/arm](http://get.maxcdn.com/maxcurl/freebsd/arm/maxcurl) ([md5](http://get.maxcdn.com/maxcurl/freebsd/arm/maxcurl.md5))
    - [windows/386](http://get.maxcdn.com/maxcurl/windows/386/maxcurl.exe) ([md5](http://get.maxcdn.com/maxcurl/windows/386/maxcurl.exe.md5))
    - [windows/amd64](http://get.maxcdn.com/maxcurl/windows/amd64/maxcurl.exe) ([md5](http://get.maxcdn.com/maxcurl/windows/amd64/maxcurl.exe.md5))
- maxtail
    - [linux/386](http://get.maxcdn.com/maxtail/linux/386/maxtail) ([md5](http://get.maxcdn.com/maxtail/linux/386/maxtail.md5))
    - [linux/amd64](http://get.maxcdn.com/maxtail/linux/amd64/maxtail) ([md5](http://get.maxcdn.com/maxtail/linux/amd64/maxtail.md5))
    - [linux/arm](http://get.maxcdn.com/maxtail/linux/arm/maxtail) ([md5](http://get.maxcdn.com/maxtail/linux/arm/maxtail.md5))
    - [darwin/386](http://get.maxcdn.com/maxtail/darwin/386/maxtail) ([md5](http://get.maxcdn.com/maxtail/darwin/386/maxtail.md5))
    - [darwin/amd64](http://get.maxcdn.com/maxtail/darwin/amd64/maxtail) ([md5](http://get.maxcdn.com/maxtail/darwin/amd64/maxtail.md5))
    - [freebsd/386](http://get.maxcdn.com/maxtail/freebsd/386/maxtail) ([md5](http://get.maxcdn.com/maxtail/freebsd/386/maxtail.md5))
    - [freebsd/amd64](http://get.maxcdn.com/maxtail/freebsd/amd64/maxtail) ([md5](http://get.maxcdn.com/maxtail/freebsd/amd64/maxtail.md5))
    - [freebsd/arm](http://get.maxcdn.com/maxtail/freebsd/arm/maxtail) ([md5](http://get.maxcdn.com/maxtail/freebsd/arm/maxtail.md5))
    - [windows/386](http://get.maxcdn.com/maxtail/windows/386/maxtail.exe) ([md5](http://get.maxcdn.com/maxtail/windows/386/maxtail.exe.md5))
    - [windows/amd64](http://get.maxcdn.com/maxtail/windows/amd64/maxtail.exe) ([md5](http://get.maxcdn.com/maxtail/windows/amd64/maxtail.exe.md5))

> Note: As of yet, these binaries have not been tested on all OS/ARCH combinations.

