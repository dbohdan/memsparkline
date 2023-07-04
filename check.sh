#! /bin/sh
set -eu

poetry run black memsparkline
poetry run mypy memsparkline
