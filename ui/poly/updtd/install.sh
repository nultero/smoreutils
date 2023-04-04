#!/usr/bin/env bash

xname="updtd"
go build -trimpath -ldflags="-s -w" .
mv $xname ~/polybar
printf "\x1b[32minstalled $xname\x1b[0m"