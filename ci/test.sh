#! /bin/sh

cd "$(dirname "$0")"/../ || exit 1

status=0
printf '=== mypy\n'
poetry run mypy memsparkline || status=1
printf '=== Tests\n'
python3 tests/test_memsparkline.py || status=1
exit "$status"
