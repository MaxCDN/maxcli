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

