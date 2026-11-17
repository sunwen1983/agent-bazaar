#!/usr/bin/env bash
# Regenerate the three static origin bundles from networks/*.hocon.
# Requires the pinned nsdeploy checkout as a sibling directory (../nsdeploy,
# commit 4e4cb36) with its venv installed, or NSDEPLOY_VENV pointing at a
# venv with nsdeploy installed. Writes into public/ (committed: Pages
# deploys that directory directly, no build step).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
NSDEPLOY_DIR="${NSDEPLOY_DIR:-$ROOT/../nsdeploy}"

if [ -n "${NSDEPLOY_VENV:-}" ]; then
  PY="$NSDEPLOY_VENV/bin/python"
elif [ -x "$NSDEPLOY_DIR/.venv/bin/python" ]; then
  PY="$NSDEPLOY_DIR/.venv/bin/python"
else
  echo "nsdeploy venv not found. Set NSDEPLOY_VENV or place the pinned" >&2
  echo "nsdeploy checkout (commit 4e4cb36) at $NSDEPLOY_DIR with .venv." >&2
  exit 1
fi

# nsdeploy is run from its own checkout (not pip-installed into the venv),
# so it needs to be on the module path.
export PYTHONPATH="$NSDEPLOY_DIR${PYTHONPATH:+:$PYTHONPATH}"

# Planned origins (placeholders until the Cloudflare Pages projects exist).
ORIGIN_A="https://agent-bazaar.pages.dev"
ORIGIN_B="https://agent-bazaar-flights.pages.dev"
ORIGIN_C="https://agent-bazaar-stays.pages.dev"

TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

echo "== origin A: concierge (links flights/stays/local-guide)"
"$PY" -m nsdeploy "$ROOT/networks/concierge.hocon" \
  --static-bundle --out "$ROOT/public/origin-a" --site-url "$ORIGIN_A" \
  --link "flights=https://agent-bazaar-flights.pages.dev/flights" \
  --link "stays=https://agent-bazaar-stays.pages.dev/stays" \
  --link "local-guide=https://agent-bazaar-stays.pages.dev/local-guide"

echo "== origin B: flights"
"$PY" -m nsdeploy "$ROOT/networks/flights.hocon" \
  --static-bundle --out "$ROOT/public/origin-b" --site-url "$ORIGIN_B"

echo "== origin C: stays + local-guide (merged)"
"$PY" -m nsdeploy "$ROOT/networks/stays.hocon" \
  --static-bundle --out "$ROOT/public/origin-c" --site-url "$ORIGIN_C"
"$PY" -m nsdeploy "$ROOT/networks/local-guide.hocon" \
  --static-bundle --out "$TMP/guide" --site-url "$ORIGIN_C"
"$ROOT/scripts/merge-origin-c.py" "$ROOT/public/origin-c" "$TMP/guide"

echo "== bundles regenerated under $ROOT/public"
find "$ROOT/public" -maxdepth 2 -name "agent-manifest.json" -o -maxdepth 2 -name "*.network.json" | sort
