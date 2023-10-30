#!/usr/bin/env sh
set -e

echo "migrate app"
./migrate db init
./migrate db migrate
echo "start app"
./server
exec "$@"
