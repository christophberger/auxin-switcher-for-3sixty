#!/usr/bin/env bash
GOOS=linux GOARCH=arm go build && scp 3sixty pi@raspi.local:/home/pi/