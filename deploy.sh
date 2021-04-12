#!/bin/sh
export PATH=$PATH:/usr/local/go/bin

# Building our app
GOOS=linux GOARCH=amd64 go build
