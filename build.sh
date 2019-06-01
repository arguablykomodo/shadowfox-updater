#!/usr/bin/env bash

OSLIST=(windows darwin linux)
ARCHLIST=(386 amd64)

OSNAMES=(windows mac linux)
ARCHNAMES=(x32 x64)

for I in ${!OSLIST[*]}; do
  for J in ${!ARCHLIST[*]}; do
    export GOOS=${OSLIST[$I]}
    export GOARCH=${ARCHLIST[$J]}
    OUT=dist/shadowfox_${OSNAMES[$I]}_${ARCHNAMES[$J]}
    if [ $GOOS = windows ]; then
      OUT=${OUT}.exe
    fi
    
    go build -o $OUT -ldflags "-X main.tag=$(git describe)"
  done
done
