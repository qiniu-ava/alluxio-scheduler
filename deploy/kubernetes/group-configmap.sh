#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

group=$1

declare -a allGroups=("default"
"alg-pro"
"blademaster"
"general-reg"
"group-ava"
"ocr"
"terror"
"video"
"video-det"
)

getGroupIndex() {
  group=$1
  for i in "${!allGroups[@]}"; do
    if [[ "${allGroups[$i]}" = "${group}" ]]; then
        echo "${i}"
        return
    fi
  done
  echo -1
}

applyConfigmap() {
  group=$1
  index=$2
  workerPort=$((20000+$index*200))
  workerDataPort=$((20000+$index*200+1))
  workerWebPort=$((20000+$index*200+2))

  if [[ $group = "default" ]]; then
    sed "s|-\<group_name\>|-${group}|g" configmap.template | \
      sed "s|\/\<group_name\>|\/alluxio-ro|g" | \
      sed "s|\<ALLUXIO_WORKER_PORT\>|${workerPort}|g" | \
      sed "s|\<ALLUXIO_WORKER_DATA_PORT\>|${workerDataPort}|g" | \
      sed "s|\<ALLUXIO_WORKER_WEB_PORT\>|${workerWebPort}|g" > $DIR/.tmp/configmap-$group.yml
  else
    sed "s|\<group_name\>|${group}|g" configmap.template | \
      sed "s|\<ALLUXIO_WORKER_PORT\>|${workerPort}|g" | \
      sed "s|\<ALLUXIO_WORKER_DATA_PORT\>|${workerDataPort}|g" | \
      sed "s|\<ALLUXIO_WORKER_WEB_PORT\>|${workerWebPort}|g" > $DIR/.tmp/configmap-$group.yml
  fi
  kubectl apply -f $DIR/.tmp/configmap-$group.yml
}

case $group in
  all)
    for i in "${!allGroups[@]}"; do
      applyConfigmap ${allGroups[$i]} $i
    done
  ;;
  *)
    index=$(getGroupIndex $group)
    if [[ $index -ne -1 ]]; then
      applyConfigmap $group $index
    else
      echo "invalid group $group"
    fi
  ;;
esac
