#!/bin/bash

# Define the target architectures and operating systems
TARGETS="windows_amd64 darwin_arm64"

VERSION="0.1.3"

# Function to build for a specific target
build_for_target() {
  os=$1
  arch=$2
  if [ $os == "windows" ]; then
      name=tmp/trends_${VERSION}_${os}_${arch}.exe
      else
      name=tmp/trends_${VERSION}_${os}_${arch}
  fi

  GOOS=$os GOARCH=$arch templ generate .
  GOOS=$os GOARCH=$arch go build -o $name -v cmd/main.go
}

# Build for each target
for target in $TARGETS; do
  os=$(echo "$target" | cut -d'_' -f1)
  arch=$(echo "$target" | cut -d'_' -f2)
  build_for_target $os $arch
done