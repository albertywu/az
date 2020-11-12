#!/usr/bin/env bash

set -euo pipefail

echo "this should fail"
echo "now"
false
echo "this should never run"