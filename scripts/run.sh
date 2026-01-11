#!/bin/bash
set -e

# Pass all arguments to the application
go run cmd/main.go "$@"
