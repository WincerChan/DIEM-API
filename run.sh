#! /bin/bash

export GOROOT=/usr/lib/golang
export GOPATH=/root/go

(/usr/bin/go run *.go &)
sleep 1
(killall go)
