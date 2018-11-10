call setGOEnv

rm *.exe
rm *.exe~
go build src/file2consul.go
go build src/file2consul-dumb.go
