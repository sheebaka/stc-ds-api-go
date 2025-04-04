#!/bin/bash

set -x

go "build" -C ./api/.

if [[ $? -ne 0 ]]; then
  exit 1
fi

./api/api