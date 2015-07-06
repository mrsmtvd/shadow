#!/bin/bash

PACKAGE_GO_IMPORT=`go list -e -f '{{.ImportComment}}' 2>/dev/null || true`
MAIN_PACKAGE_PATH=$GOPATH"/src/"$PACKAGE_GO_IMPORT
PACKAGE_COMPRESS="true"

export GOOS=linux
export CGO_ENABLED=0

mkdir -p `dirname $MAIN_PACKAGE_PATH`
ln -sf /src $MAIN_PACKAGE_PATH

if [ -e "$MAIN_PACKAGE_PATH/Godeps/_workspace" ]; then
    GOPATH=$MAIN_PACKAGE_PATH/Godeps/_workspace:$GOPATH
else
    go get -t -d -v ./...
fi

for PACKAGE in $(go list -e -f '{{.ImportComment}}' ./... 2>/dev/null || true)
do
    echo "Build package $PACKAGE"

    PACKAGE_NAME=${PACKAGE##*/}

    cd $GOPATH"/src/"$PACKAGE

    if [ -e "./Makefile" ]; then
        make build
    else
        go build -a -v -race
    fi

    if [ $? -eq 0 ]; then
        if [[ $PACKAGE_COMPRESS == "true" ]]; then
          goupx $PACKAGE_NAME
        fi

        if [ -e "./Makefile" ]; then
            make build-post
        fi

        echo "Build package $PACKAGE SUCCESS"

        if [ -e "/var/run/docker.sock" ] && [ -e "./Dockerfile" ]; then
          docker build -t "$PACKAGE_NAME:latest" ./
        fi
    else
        echo "Build package $PACKAGE FAILED"
    fi
done