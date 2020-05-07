#!/bin/sh
export GOPATH=$PWD

go build src/file2consul.go
go build src/consul2file.go
go build src/file2consul-dumb.go

