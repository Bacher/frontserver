#!/usr/bin/env sh

go build -o build/frontserver ./server && docker build -t frontserver .