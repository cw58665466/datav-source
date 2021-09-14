#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o datav main.go
echo "build finished"

scp ./datav root@39.99.191.148:/mnt/datavsource/datav_new
echo "copy finished"

ssh root@39.99.191.148 "/mnt/datavsource/codemonitor.sh"
echo "code updated"