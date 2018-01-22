#!/usr/bin/env bash

export PROJECT_DIR="$(gb env GB_PROJECT_DIR)"
export GOPATH="$PROJECT_DIR/vendor:$PROJECT_DIR"
export PATH="$PROJECT_DIR/bin:$PROJECT_DIR/deploy:$HOME/gopath/bin:$PATH"
