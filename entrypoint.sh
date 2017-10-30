#!/usr/bin/env sh


set -e
echo $@
go run main.go
exec "$@"