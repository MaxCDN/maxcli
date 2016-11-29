maxcdn tools
============

Collection of CLI tools for interfacing with MaxCDN's REST API

> Built using [github.com/MaxCDN/go-maxcdn](https://github.com/MaxCDN/go-maxcdn).

#### Note

`maxtail` currently isn't working as intended. See [issue #2](https://github.com/MaxCDN/maxcli/issues/2) for updates.

Configuration
-------------

All tools use a configuration file as it's last means of getting `alias`, `secret` and
`token`. See individal tool `help` for addtional options available in your configuration.

```yaml
---
alias: YOUR_ALIAS
token: YOUR_TOKEN
secret: YOUR_SECRET
zones:
  - YOUR_ZONE
```

See [sample.maxcdn.yml](sample.maxcdn.yml) for a more complete example.


Installing:
-----------

This can also be installed for system wide use if your `GOBIN` is set via the following:

```bash
# via 'go get' && 'go install'
go get -u github.com/MaxCDN/maxcli/{{tool}}

# manually
##

git clone https://github.com/MaxCDN/maxcli
cd maxcli

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

A set of binaries for all tools have been prebuilt using golang's cross compiler.

**Building All Binaries:**

```bash
make build/all

# or
make build/all/{{tool}}
```

Here's what's available for each tool:

- maxreport
    - [linux-amd64](_builds/maxreport/linux/amd64/maxreport) ([md5](_builds/maxreport/linux/amd64/maxreport.md5))
    - [darwin-amd64](_builds/maxreport/darwin/amd64/maxreport) ([md5](_builds/maxreport/darwin/amd64/maxreport.md5))
    - [windows-amd64](_builds/maxreport/windows/amd64/maxreport.exe) ([md5](_builds/maxreport/windows/amd64/maxreport.exe.md5))
- maxpurge
    - [linux-amd64](_builds/maxpurge/linux/amd64/maxpurge) ([md5](_builds/maxpurge/linux/amd64/maxpurge.md5))
    - [darwin-amd64](_builds/maxpurge/darwin/amd64/maxpurge) ([md5](_builds/maxpurge/darwin/amd64/maxpurge.md5))
    - [windows-amd64](_builds/maxpurge/windows/amd64/maxpurge.exe) ([md5](_builds/maxpurge/windows/amd64/maxpurge.exe.md5))
- maxcurl
    - [linux-amd64](_builds/maxcurl/linux/amd64/maxcurl) ([md5](_builds/maxcurl/linux/amd64/maxcurl.md5))
    - [darwin-amd64](_builds/maxcurl/darwin/amd64/maxcurl) ([md5](_builds/maxcurl/darwin/amd64/maxcurl.md5))
    - [windows-amd64](_builds/maxcurl/windows/amd64/maxcurl.exe) ([md5](_builds/maxcurl/windows/amd64/maxcurl.exe.md5))
- maxtail
    - [linux-amd64](_builds/maxtail/linux/amd64/maxtail) ([md5](_builds/maxtail/linux/amd64/maxtail.md5))
    - [darwin-amd64](_builds/maxtail/darwin/amd64/maxtail) ([md5](_builds/maxtail/darwin/amd64/maxtail.md5))
    - [windows-amd64](_builds/maxtail/windows/amd64/maxtail.exe) ([md5](_builds/maxtail/windows/amd64/maxtail.exe.md5))

To cross compile your own binary for a different OS / ARCH, run the following...

```
env GOOS={{OS}} GOARCH={{ARCH}} go build github.com/MaxCDN/maxcli/{{tool}}
```

**Requires Go 1.5 or higher**

