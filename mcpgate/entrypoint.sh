#!/bin/sh
# Container entrypoint: inject the demo API key from the DEMO_API_KEY
# environment variable (a fly secret) into the gateway config, then run
# mcpgate. The key never lives in the image or the repo.
set -e

if [ -z "${DEMO_API_KEY:-}" ]; then
  echo "DEMO_API_KEY is not set — set it as a fly secret (fly secrets set DEMO_API_KEY=...)" >&2
  exit 1
fi

# The placeholder in config.yaml.template is the literal string
# "demo-key-CHANGE-ME". The key must not contain '|' or '&' (sed).
sed "s|demo-key-CHANGE-ME|${DEMO_API_KEY}|g" /app/config.yaml.template > /app/config.yaml
chmod 600 /app/config.yaml

exec /app/mcpgate -config /app/config.yaml
