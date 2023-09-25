#!/usr/bin/env sh
set -e

echo "start app"
./server
exec "$@"
