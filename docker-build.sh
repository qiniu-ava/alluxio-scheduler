#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

tag=$(git rev-parse --short=7 HEAD)

docker build \
  -t reg-xs.qiniu.io/atlab/alluxio-schedular:"$tag" \
  -f "$DIR"/app/docker/Dockerfile \
  "$DIR"/app/src
