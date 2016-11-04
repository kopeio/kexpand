#!/bin/bash

mkdir -p /go
export GOPATH=/go

mkdir -p /go/src/github.com/kopeio
ln -s /src/ /go/src/github.com/kopeio/kexpand

cd /go/src/github.com/kopeio/kexpand
if [ ! -d vendor ]; then
  go get -t -d -v ./...;
fi
make gocode_docker

mkdir -p /src/.build/artifacts/
cp /go/bin/kexpand /src/.build/artifacts/
