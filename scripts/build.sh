#!/bin/bash



echo $(pwd)
rm -rf ./dist/
mkdir ./dist/

APP_NAME=$1
TAG=$2

if [ "X${APP_NAME}" = "X" ];then
    echo "app name cannot be empty"
    exit 1
fi

if [ "X${TAG}" = "X" ];then
    echo "tag cannot be empty, example: ./build.sh appname v1.5.8"
    exit 1
fi



COMMIT=`git rev-parse HEAD`
BUILDTIME=`date +'%Y-%m-%d_%T'`
GOVER=`go env GOVERSION`
UTILS="github.com/qvcloud/gopkg/version"

FLAGS="-s -w -X ${UTILS}.Version=${TAG} -X ${UTILS}.Commit=${COMMIT} -X ${UTILS}.Build=${BUILDTIME} -X ${UTILS}.Go=${GOVER}"

echo ${FLAGS}

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="${FLAGS}" -o ./dist/${APP_NAME} cmd/main.go
