#!/bin/bash

export ENVCONF_LOAD_DOTFILE=disable
export DB_USER=test
export DB_HOST=localhost
export DB_PORT=5432

go run main.go
