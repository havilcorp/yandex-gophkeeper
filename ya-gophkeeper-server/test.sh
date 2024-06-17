#!/bin/sh

go test ./... -coverprofile cover.tmp.out
cat cover.tmp.out | grep -v "/mocks/" > cover.out
go tool cover -func=cover.out | grep total:
go tool cover -html=cover.out