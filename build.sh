#!/bin/bash

GOOS=linux GOARCH=mipsle go build -o hamgo_mipsel
GOOS=windows GOARCH=amd64 go build -o hamgo_win.exe
GOOS=linux GOARCH=amd64 go build -o hamgo_x86
