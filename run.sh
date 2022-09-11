#!/usr/bin/sh

go build
./GuthiNetwork -port 7000 &
./GuthiNetwork