#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$SCRIPT_DIR/.."

echo "Generating OpenAPI spec..."
cd "$ROOT_DIR/trails"
go run ./cmd/app --generate-openapi

echo "Generating TypeScript types..."
cd "$ROOT_DIR/ui"
npm run gen:api

echo "Done."
