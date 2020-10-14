#! /bin/sh
set -e

cd "$(dirname "$0")"
../memsparkline sleep 1 2>&1 | grep max > /dev/null
