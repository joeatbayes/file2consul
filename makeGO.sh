#!/bin/sh
export GOPATH=$PWD

go build src/GenericHTTPTestClient.go
go build src/File2Console.go
go build src/ex-command-args.com
go build src/classifyFiles.go

