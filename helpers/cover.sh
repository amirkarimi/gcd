#!/bin/sh
PACKAGES=$(go list ./... | grep -v /vendor/)
echo "mode: count" > coverage-all.out
for pkg in $PACKAGES
do
  go test -cover -coverprofile=coverage.out -covermode=count $pkg
  tail -n +2 coverage.out >> ./coverage-all.out
done
go tool cover -html=coverage-all.out
