#!/bin/bash

mkdir -p /go
export GOPATH=/go

mkdir -p /go/src/github.com/kopeio
ln -s /src/ /go/src/github.com/kopeio/kexpand

cd /go/src/github.com/kopeio/kexpand
make gocode

mkdir -p /src/.build/artifacts/
cp /go/bin/kexpand /src/.build/artifacts/
