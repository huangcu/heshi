#!/usr/bin/env bash
. env_windows.bash


go build -o bin/heshi_service -race heshi_service
go build -o bin/cmd -race heshi_service/cmd

# export STAGE=dev
# export TRACE=true

heshi_service
