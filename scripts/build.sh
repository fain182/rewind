#!/bin/sh
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR/..
GOOS=linux GOARCH=amd64 go build rewind.go
rm -r dist
mkdir -p dist
mv rewind dist/
cp -r assets dist/
