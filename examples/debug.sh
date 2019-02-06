#!/bin/bash
echo "Building application"
export cur=$(pwd)
export PATH="$PATH:$cur"
go build -gcflags "all=-N -l" sigkeys_example.go

echo "Starting the application with dlv for debugging, connect to localhost:2345"
dlv --listen=127.0.0.1:2345 --headless=true --api-version=2 exec sigkeys_example -- "$@"
