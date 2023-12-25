#!/bin/sh

if [ -f main ]; then
    rm -rf main
fi

go build main.go
