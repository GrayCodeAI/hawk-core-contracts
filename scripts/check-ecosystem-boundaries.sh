#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

# hawk-core-contracts is a foundation library: it sits below every engine
# and below hawk itself. It must depend on nothing in the hawk ecosystem —
# only the Go standard library. Imports of this module's own path (internal
# cross-package imports between types/, tools/, events/, etc.) are fine.
FORBIDDEN='github\.com/GrayCodeAI/(?!hawk-core-contracts(/|"))'
FORBIDDEN_GREP='github\.com/GrayCodeAI/[^"]*'
SELF_MODULE='github.com/GrayCodeAI/hawk-core-contracts'

if command -v rg >/dev/null 2>&1; then
  violations="$(rg -n "$FORBIDDEN" --pcre2 --glob '*.go' . || true)"
else
  violations="$(grep -rn --include='*.go' -E "$FORBIDDEN_GREP" . | grep -v "$SELF_MODULE" || true)"
fi

if [[ -n "${violations}" ]]; then
  echo "forbidden hawk-eco imports found in hawk-core-contracts:"
  echo "${violations}"
  echo
  echo "hawk-core-contracts is a foundation repo — it must not depend on hawk, engines, or any other GrayCodeAI/* package"
  exit 1
fi

echo "ecosystem boundary guard passed"
