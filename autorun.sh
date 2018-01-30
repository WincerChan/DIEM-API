#! /bin/bash

export GOROOT=/usr/lib/golang
export GOPATH=/root/go
echo "GOPATH:$GOPATH"


pwd
/usr/bin/go run *.go &
killall go
