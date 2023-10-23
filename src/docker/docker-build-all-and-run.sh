#!/bin/bash

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [[ "$PWD" != *"$script_dir"* ]]; then
  cd $script_dir
fi

cd ../../../institute-mongodb
./docker-build.sh

cd ../institute-person-api/src/docker
./docker-build.sh

docker image prune -f

docker compose up --detach
