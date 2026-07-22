#!/usr/bin/env bash
# Release tagging script for hawk-core-contracts.
# Ensures git status is clean, runs checks, updates VERSION, and creates a semver tag.
set -euo pipefail

if [ $# -ne 1 ]; then
  echo "Usage: $0 <version>"
  echo "Example: $0 v0.1.0"
  exit 1
fi

VERSION="$1"

if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.]+)?$ ]]; then
  echo "Error: Version must follow semver format (e.g. v0.1.0 or v1.0.0-rc.1)"
  exit 1
fi

if [ -n "$(git status --porcelain)" ]; then
  echo "Error: Working directory is dirty. Commit or stash changes before tagging."
  exit 1
fi

echo "Running pre-tag release checks..."
make test || exit 1

# Strip leading 'v' for VERSION file
RAW_VERSION="${VERSION#v}"
echo "$RAW_VERSION" > VERSION

if [ -f "version.go" ]; then
  sed -i '' "s/Version = \".*\"/Version = \"$RAW_VERSION\"/" version.go || true
fi

git add VERSION version.go 2>/dev/null || git add VERSION
git commit -m "chore(release): prepare $VERSION"
git tag -a "$VERSION" -m "Release $VERSION"

echo "Successfully tagged $VERSION!"
echo "Push release with: git push origin main --tags"
