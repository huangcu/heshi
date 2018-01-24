#!/usr/bin/env bash
. env_windows.bash

go build -o bin/heshi_service -tags dev heshi_service
go build -o bin/cmd -tags dev cmd

# export STAGE=dev
# export TRACE=true

heshi_service
