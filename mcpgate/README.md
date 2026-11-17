# mcpgate deployment for Agent Bazaar

mcpgate (the author's own MCP gateway, `sunwen1983/mcpgate`) fronts the
three mock MCP backends in `../tools/mock-mcp`. Browsers call it
cross-origin from the three Pages sites; the CORS middleware merged in
mcpgate PR #20 (commit `9bb4a88`) makes that possible.

## Files

| File | Purpose |
|------|---------|
| `config.yaml` | Gateway config: 1 demo API key (placeholder value — **replace before any public deploy**), 3 stdio backends, cache + breaker settings |
| `LOCAL_E2E.md` | Local verification: commands + results (all passed 2026-11-17) |
| `Dockerfile` | Multi-stage image: builds mcpgate (pinned `9bb4a88`) + mock-mcp, runs the gateway |
| `entrypoint.sh` | Substitutes the demo key from the `DEMO_API_KEY` env (fly secret) into the config at container start, then execs the gateway |
| `fly.toml` | fly.io app `agent-bazaar-mcpgate` (region `dfw`, 1 shared-cpu-1x 256 MB machine, always on) |

## Run locally

```bash
# from the agent-bazaar repo root
(cd tools/mock-mcp && go build -o mock-mcp .)
go build -o /tmp/mcpgate ~/workspace/mcpgate/cmd/mcpgate
/tmp/mcpgate -config mcpgate/config.yaml
```

Listens on `:8080`. See `LOCAL_E2E.md` for the curl checks.

## Deploy (Wen — see `../docs/deploy.md` for the full checklist)

```bash
cd mcpgate
fly launch --no-deploy          # first time only; accept the generated app name check
fly secrets set DEMO_API_KEY="$(openssl rand -hex 16)"
fly deploy                      # builds with the remote builder, no local Docker needed
curl https://agent-bazaar-mcpgate.fly.dev/healthz
```

## Rate limit, by design

The demo key is public (it ships in the bazaar site's frontend config),
so the gateway is the only defense: 30 req/min sustained (0.5/s,
burst 10). Burst 10 lets one full team run (~6–10 tool calls, spread
over LLM round-trips) through; sustained abuse is capped. The backends
are read-only mock data, so there is nothing else to abuse. Rotate the
key after the contest.
