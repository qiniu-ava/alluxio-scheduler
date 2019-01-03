#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

kubectl create configmap alluxio-schedular-config --from-file=nodelist --from-file=grouplist
