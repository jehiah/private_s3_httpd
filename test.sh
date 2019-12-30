#!/bin/bash
set -e

go test -timeout 60s ./...
# gb test -timeout 60s -race
mkdir -p bin
go build -o bin/private_s3_httpd ./src/cmd/private_s3_httpd
