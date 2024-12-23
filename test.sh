#!/bin/sh

function cleanup() {
    rm -rf *.jpg
}

echo "Test crop single image..."
go run main.go crop --src kodak/00001.jpg
count=`ls *.jpg | wc -l`
[[ $count = 2 ]] || exit 1
cleanup

echo "Test crop images inside directory..."
mkdir dist
go run main.go crop --src kodak --dest dist
count=`ls dist/*.jpg | wc -l`
[[ $count = 4 ]] || exit 1
rm -rf dist
cleanup

