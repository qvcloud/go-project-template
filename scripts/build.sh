#!/bin/bash



echo $(pwd)
rm -rf ./dist/
mkdir ./dist/

TAG=$1
if [ "X${TAG}" = "X" ];then
    echo "tag cannot be empty, example: ./build.sh v1.5.8"
    exit 1
fi



COMMIT=`git rev-parse HEAD`
BUILDTIME=`date +'%Y-%m-%d_%T'`
GOVER=`go env GOVERSION`
UTILS="github.com/qvcloud/gopkg/version"

FLAGS="-s -w -X ${UTILS}.Version=${TAG} -X ${UTILS}.Commit=${COMMIT} -X ${UTILS}.Build=${BUILDTIME} -X ${UTILS}.Go=${GOVER}"

echo ${FLAGS}

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="${FLAGS}" -o ./dist/xxName cmd/main.go 
