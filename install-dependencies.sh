#!/usr/bin/env bash

set -xe

mkdir -p $GOPATH/src/github.com/gogo
git clone https://github.com/gogo/protobuf $GOPATH/src/github.com/gogo/protobuf
mkdir -p $GOPATH/src/github.com/katydid
git clone https://github.com/katydid/katydid $GOPATH/src/github.com/katydid/katydid
