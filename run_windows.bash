#!/usr/bin/env bash
. env_windows.bash

go build -o bin/heshi_service -tags dev heshi_service

# export STAGE=dev
# export TRACE=true

heshi_service
