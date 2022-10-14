#!/bin/env bash

all: sn

sn: main.go
	go build -o sn main.go
	cp sn /usr/local/bin/
