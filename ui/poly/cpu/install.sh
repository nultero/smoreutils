#!/bin/bash

go build . &&
mv ucpu ~/polybar &&
polybar-msg cmd restart