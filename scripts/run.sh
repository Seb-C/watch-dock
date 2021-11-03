#!/bin/bash

set -e

./scripts/build.sh

export WDBIN="$(realpath ./bin/watch-dock)"
cd ./example/hello-world.go
"$WDBIN" -- -f docker-compose.yaml
