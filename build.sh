#!/bin/bash

set -e

echo "Building soloterm for all platforms..."

CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -o bin/soloterm_Darwin_arm64 . && echo "  macOS (arm64)"
CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -o bin/soloterm_Darwin_x86_64 . && echo "  macOS (x86_64)"
CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build -o bin/soloterm_Linux_arm64 . && echo "  Linux (arm64)"
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -o bin/soloterm_Linux_x86_64 . && echo "  Linux (x86_64)"
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o bin/soloterm_Windows_arm64.exe . && echo "  Windows (arm64)"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/soloterm_Windows_x86_64.exe . && echo "  Windows (x86_64)"

echo "Done."
