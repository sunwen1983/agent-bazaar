# Local end-to-end verification (HTTP layer)

Run: 2026-11-17, VM. Gateway built from `~/workspace/mcpgate` main
(`9bb4a88`, includes the CORS middleware from PR #20). Mock backends
built from `tools/mock-mcp` in this repo.

## Setup

```bash
(cd tools/mock-mcp && go build -o mock-mcp .)
go build -o /tmp/mcpgate-e2e/mcpgate ~/workspace/mcpgate/cmd/mcpgate
/tmp/mcpgate-e2e/mcpgate -config mcpgate/config.yaml   # from repo root
```

Gateway log on start (abridged):

```
{"level":"INFO","msg":"mcpgate listening","addr":":8080","backends":3,"api_keys":1}
mock-mcp flights: ready (deterministic sample data)
mock-mcp stays: ready (deterministic sample data)
mock-mcp activities: ready (deterministic sample data)
```

## Checks (all passed)

`KEY="demo-key-CHANGE-ME"` — the placeholder value in `config.yaml`,
used as the test key. It is not a secret; it must be replaced with a
generated key before any public deployment.

| # | Check | Command (abridged) | Result |
|---|-------|--------------------|--------|
| 1 | No key → 401 | `POST /mcp/flights tools/list`, no key | `{"error":{"code":-32001,"message":"missing or invalid API key"}}` HTTP 401 |
| 2 | Wrong key → 401 | same, `X-API-Key: wrong` | HTTP 401 |
| 3 | tools/list × 3 | `POST /mcp/{flights,stays,activities}` with key | `flights → ['search_flights']`, `stays → ['search_stays']`, `activities → ['search_activities']` |
| 4 | tools/call cache MISS → HIT | `search_flights {adults: 2}` twice | 1st `X-Cache: MISS`, 2nd `X-Cache: HIT`; cheapest Southwest $356, 6 options |
| 5 | tools/call stays + activities | `search_stays {nights: 2}`, `search_activities {interests: ["live music","food"]}` | stays cheapest Austin Motel $236 (8 options); activities matched 5 |
| 6 | OPTIONS preflight | `OPTIONS /mcp/flights`, `Origin: https://example.com` | HTTP 200, `Access-Control-Allow-Origin: *`, `Allow-Methods: GET, POST, OPTIONS`, `Allow-Headers: X-API-Key, Authorization, Content-Type`, `Max-Age: 86400` |
| 7 | CORS headers on real POST | `POST /mcp/flights` with `Origin:` | Same `Access-Control-Allow-*` headers present |
| 8 | Unknown backend → 404 | `POST /mcp/nope` | `{"error":{"code":-32601,"message":"unknown backend: nope"}}` HTTP 404 |
| 9 | /metrics stays same-origin | `OPTIONS /metrics` with `Origin:` | No `Access-Control-*` headers (as intended) |
| 10 | Rate limit (30/min) | 15 rapid `tools/list` after a 20 s idle | 10× HTTP 200, then 5× HTTP 429 (burst 10, then 0.5/s refill) |
| 11 | Budget | cheapest combo | flights $356 + stays $236 + free activities $0 = **$592 < $800** |

Gateway log had 0 error lines during the run. Gateway stopped
afterwards; port 8080 free.

## What this does NOT verify

Browser execution (consent popup, cross-origin peer loading, the
nsDeploy runtime calling this gateway) — the VM has no working browser
for that (see MEMORY.md: cloud browser can't reach VM localhost,
local Chromium blocked by policy, tunnels blocked). That check happens
after the Cloudflare Pages deploy, against the public URLs. It stays
with Mila.
