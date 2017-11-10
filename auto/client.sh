#!/bin/bash

wget --no-verbose https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz
sudo tar -xzf $HOME/go1.8.3.linux-amd64.tar.gz -C /usr/local/
sudo ln -s /usr/local/go/bin/go /usr/local/bin/go
go version

cd $HOME/src/github.com/vasili-v/grpc-stream-test/gst-client
GOPATH=$HOME bash -c 'go install -v'
