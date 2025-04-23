#!/bin/bash
set -e

echo "Installing optique"

go build -o optique main.go
mv optique ~/.local/bin
