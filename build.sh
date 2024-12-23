#!/bin/bash

oss=("linux" "darwin")
archs=("amd64" "arm64")
for os in ${oss[@]}; do
    for arch in ${archs[@]}; do
        GOOS=$os GOARCH=$arch go build -o hf-utils_$os-$arch main.go
    done
done
