# mock-mcp — deterministic MCP backends for Agent Bazaar

One Go program, three datasets. Each process is a minimal MCP server
over stdio (JSON-RPC 2.0, one object per line), modeled on mcpgate's
`examples/echo-backend`. mcpgate spawns one process per backend.

## Build

```bash
go build -o mock-mcp .
```

(Go 1.23+; the VM toolchain is `~/go-install/go/bin/go`.)

## Run

```bash
./mock-mcp flights      # search_flights
./mock-mcp stays        # search_stays
./mock-mcp activities   # search_activities
```

Or through mcpgate (from the agent-bazaar repo root):

```bash
# build mcpgate from the local clone (main includes the CORS middleware, PR #20)
go build -o /tmp/mcpgate ~/workspace/mcpgate/cmd/mcpgate
/tmp/mcpgate -config mcpgate/config.yaml
```

Then `curl` the gateway — see `../mcpgate/LOCAL_E2E.md` for the
exact checks (tools/list, tools/call, cache MISS→HIT, 401, preflight).

## Data

Fixed sample inventory, not live availability. Every response says so
in its `note` field. Same arguments always return the same answer.

- **flights** — 6 round-trip SFO→AUS options, $178–$268 per person.
  Cheapest for 2 adults: $356.
- **stays** — 8 Austin hotels, $118–$210/night. Cheapest 2 nights: $236.
- **activities** — 10 Austin activities, $0–$30 per person, three free.

Cheapest full weekend (flights $356 + hotel $236 + free activities)
is $592; a comfortable mid-range pick (Alaska $392 + Carpenter $298 +
~$50 of activities) is ~$740 — both under the demo's $800 budget.

Tool arguments are filters only (max price, neighborhood, interests,
day); dates are echoed back when provided. Unknown tool names get
JSON-RPC `-32602`; unknown methods get `-32601`.

## Swap in real APIs later

The tool names, schemas, and response shapes are the contract the
HOCON networks (`networks/*.hocon`) rely on. To replace the mock with a
real provider, keep the same tool names and response fields and point
mcpgate at the real MCP server instead — no network change needed.
