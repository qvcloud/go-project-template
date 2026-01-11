#!/bin/bash

dlv debug cmd/main.go --headless --listen=:2345 --api-version=2 --log