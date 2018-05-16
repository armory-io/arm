#!/bin/bash
set -e

VERSION=${VERSION:-"HEAD"}

# compile enabling the version check and cli version
ldflags="-X github.com/armory-io/arm/cmd.enableVersionCheck=check -X github.com/armory-io/arm/cmd.currentVersion=${VERSION}"

echo "building linux binary"
GOOS=linux go build -o arm_linux_amd64 -ldflags $ldflags
echo "building osx binary"
GOOS=darwin go build -o arm_darwin_amd64 -ldflags $ldflags
echo "building windows binary"
GOOS=windows go build -o arm_windows_amd64.exe -ldflags $ldflags