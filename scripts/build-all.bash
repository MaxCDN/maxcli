#!/bin/bash

archs=`ls $(go env GOROOT)/pkg | grep -v "obj\|tool\|race" | grep "darwin\|freebsd\|linux\|freebsd\|openbsd\|solaris\|windows`

if ! test "$archs"; then
  echo "no valid os/arch pairs were found to build"
  exit 1
fi

tools="maxcurl maxpurge maxreport maxtail"

# allow for single tool build
test "$1" && tools="$1"

for tool in $tools
do
  src="$tool/$tool.go"
  for arch in $archs
  do
    split=(${arch//_/ })
    goos=${split[0]}
    goarch=${split[1]}

    src=$tool/$tool.go
    target=$tool/builds/$goos/$goarch/$tool

    [[ "windows" == "$goos" ]] && target=$target.exe

    if ! test -d "$(dirname $target)"
    then
      echo "mkdir -pv $(dirname $target)"
      mkdir -pv $(dirname $target)
    fi

    echo "GOOS=$goos GOARCH=$goarch go build -x -o $target $src"
    GOOS=$goos GOARCH=$goarch go build -x -o $target $src
    md5sum $target | tee $target.md5
  done
done
