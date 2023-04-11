#!/usr/bin/env bash

xname="updt"
go build -trimpath -ldflags="-s -w" .
mv $xname ~/.$USER/bin
printf "\x1b[32minstalled $xname\x1b[0m"