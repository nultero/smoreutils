#!/usr/bin/env bash

go build -trimpath -ldflags="-s -w" .
mv upstatd ~/polybar
printf "\x1b[32minstalled upstatd\x1b[0m"