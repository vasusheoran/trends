#!/bin/bash

# Define the target architectures and operating systems
TARGETS="darwin_arm64_clang windows_amd64_x86_64-w64-mingw32-gcc"

VERSION="0.1.4"

# Function to build for a specific target
build_for_target() {
  os=$1
  arch=$2
  cc=$3
  if [ $os == "windows" ]; then
      name=tmp/trends_${VERSION}_${os}_${arch}.exe
      else
      name=tmp/trends_${VERSION}_${os}_${arch}
  fi

  GOOS=$os GOARCH=$arch CC=$cc CGO_ENABLED='1' templ generate .
  GOOS=$os GOARCH=$arch CC=$cc CGO_ENABLED='1' go build -o $name -v cmd/main.go
}

# Build for each target
for target in $TARGETS; do
  os=$(echo "$target" | cut -d'_' -f1)
  arch=$(echo "$target" | cut -d'_' -f2)
  cc=$(echo "$target" | cut -d'_' -f3 -f4)
  build_for_target $os $arch $cc
done