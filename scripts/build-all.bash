#!/bin/bash

set -ue

function gobuild {
  local tool=$1
  local os=$2
  local arch=$3
  local builddir=_builds/$tool/$os/$arch
  local target=$builddir/$tool
  local src=$tool/${tool}.go

  [[ "windows" == "$os" ]] && target=${target}.exe

  echo " "
  echo "---"
  echo "building $tool for $os $arch"
  echo "      to $target"

  mkdir -p $builddir

  env GOOS=$os GOARCH=$arch go build -o $target $src
  md5sum $target | tee $target.md5
}

tools="maxcurl maxpurge maxreport maxtail"

# allow for single tool build
#test "$1" && tools="$1"


for tool in $tools
do
  gobuild $tool windows amd64
  gobuild $tool darwin amd64
  gobuild $tool linux amd64
done
