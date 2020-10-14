#! /bin/sh
set -e

cd "$(dirname "$0")"

../memsparkline sleep 1 2>&1 | grep max > /dev/null

../memsparkline -w 2000 sleep 1 2>&1 | wc -l | grep '^3$' > /dev/null
../memsparkline -n -w 100 sleep 1 2>&1 | wc -l | grep '^1[0-9]$' > /dev/null
