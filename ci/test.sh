#! /bin/sh
set -e

python3 "$(dirname "$0")/../tests/test_memsparkline.py"
