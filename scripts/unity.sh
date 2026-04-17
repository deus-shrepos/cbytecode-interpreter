#!/usr/bin/env bash

set -euo pipefail

COMMIT="v2.6.1"
BASE="https://raw.githubusercontent.com/ThrowTheSwitch/Unity/${COMMIT}/src"
DEST="vendor/unity"

mkdir -p "$DEST"

for f in unity.c unity.h unity_internals.h; do
    echo "fetching ${f}..."
    curl -fsSL "${BASE}/${f}" -o "${DEST}/${f}"
done

echo "Done."