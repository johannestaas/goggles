#!/bin/bash
set -euo pipefail
find . -type f -name '*.go' -exec go fmt {} \;
