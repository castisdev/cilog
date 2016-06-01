#!/usr/bin/env bash
DIR=$GOPATH/src/github.com/castisdev/cilog
cd $DIR
go get ./...
go build -x ./...
go get -t ./...
go get github.com/axw/gocov/...
go get github.com/Centny/gocov-xml
go get github.com/jstemmer/go-junit-report
go get github.com/golang/lint/golint
go test -v ./... | go-junit-report > junit.xml
gocov test ./... | gocov-xml -b $GOPATH > coverage.xml
golint ./... > lint.txt
go tool vet -printf=false . > vet.txt
