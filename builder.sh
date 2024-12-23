#!/bin/sh

go build -o "c" dev_scripts.go
sudo mv c "/usr/local/bin/c"
sudo chmod 750 /usr/local/bin/c
