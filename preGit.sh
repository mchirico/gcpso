#!/bin/bash

#  git push --delete origin tagName
#  git tag -d tagName
#
# git tag -a v1.0 -m 'simple main test with flags'
# git checkout tags/v1.0



gofmt -s -w examples/example.go
gofmt -s -w example/example_test.go
go test -v -coverprofile=c0.out -covermode=atomic github.com/mchirico/gcpso/hello
go test -v -coverprofile=c1.out -covermode=atomic github.com/mchirico/gcpso/utils
go vet -v github.com/mchirico/gcpso/examples github.com/mchirico/gcpso/utils


# Mock just area
# mockgen -destination=./mocks/mock_stuff.go -package=mocks github.com/mchirico/gcpso/configs/mocks Area

