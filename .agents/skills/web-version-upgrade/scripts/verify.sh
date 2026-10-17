#!/bin/sh
set -eu

repo_root=$(CDPATH= cd "$(dirname "$0")/../../../.." && pwd)
web_dir="$repo_root/apps/web"

cd "$web_dir"
npm run lint
npm run format:check
npm run build
npm ls --depth=0
