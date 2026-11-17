# Agent Bazaar — agents from different websites team up in your browser

**Open Agentic Web Challenge entry** (Devpost: `wen-sun1217`, solo).
Submission deadline: **2027-01-11 10:00 AM CT**. Prize pool: $10,000 cash.

Agents from different websites meet, verify each other, and team up — live in
your browser, with no central orchestrator. The browser is the orchestrator;
[mcpgate](https://github.com/sunwen1983/mcpgate) is the tool backend.

## What this is

A public agent website built on the
[MVP Profile 0.2](https://github.com/cognizant-ai-labs/agentic-web/blob/main/docs/worldwide_agentic_web_mvp_spec.md)
spec, using the [nsDeploy](https://github.com/cognizant-ai-lab/nsdeploy)
reference runtime (pinned commit `4e4cb36`).

- **Origin A** — the bazaar itself (`concierge` network): one-sentence trip
  request in, a full Austin weekend itinerary out.
- **Origin B** — `flights` peer network: round-trip flights to Austin.
- **Origin C** — `stays` + `local-guide` peer networks: hotels and a
  day-by-day Austin plan.
- **Tool backend** — three MCP services (flights / stays / activities) behind
  mcpgate, returning deterministic mock data (no real bookings, no payments).
- **Trust, visible** — every peer is checked before it joins the team:
  manifest reachable + name matches, Ed25519 signature verifies, local-trust
  status shown. A bad peer is refused on screen (fail closed). No NANDA
  registry lookup yet — the UI keeps the slot for it (see `docs/`).

The official nsDeploy travel example is itself a 10-agent travel concierge,
so the travel content is **not** the novelty here. The novelty is the
marketplace: cross-origin discovery, visible verification, and teaming-up in
the browser.

## Repo layout

```
networks/            HOCON network definitions (the source of truth)
  concierge.hocon    origin A front network; links flights/stays/local-guide
  flights.hocon      origin B
  stays.hocon        origin C
  local-guide.hocon  origin C
scripts/
  build-bundles.sh   regenerates public/ from networks/ via nsdeploy
  merge-origin-c.py  merges the stays + guide bundles into one origin
public/              generated static bundles (what gets deployed)
  origin-a/          bazaar main site
  origin-b/          flights peer site
  origin-c/          stays + local-guide peer site
docs/
  demo-key.md        how the demo LLM key works (placeholder only, no real key)
```

Note: the generated bundles keep nsDeploy's real on-disk layout
(`agent-manifest.json`, `<network>.network.json`, `links.json`, `runtime/` at
the bundle root). The plan's §5.6 sketch showed an `/agents/` prefix; the
actual nsDeploy layout wins.

## Build

Requires the pinned nsDeploy checkout as a sibling directory
(`../nsdeploy`, commit `4e4cb36`) with its Python venv installed, or set
`NSDEPLOY_VENV` to a venv with nsdeploy installed.

```bash
./scripts/build-bundles.sh
```

Regenerates `public/origin-a`, `public/origin-b`, `public/origin-c` and
smoke-checks them with a plain HTTP server (key resources must be 200;
CORS headers come from the generated `_headers`, enforced by the host —
Cloudflare Pages — not by the smoke server).

## Deploy (planned, not done yet)

- Static sites: three Cloudflare Pages projects (one per origin).
  Placeholder origins used in the bundles — confirm when creating the
  projects:
  - A: `https://agent-bazaar.pages.dev`
  - B: `https://agent-bazaar-flights.pages.dev`
  - C: `https://agent-bazaar-stays.pages.dev`
- Tool gateway: mcpgate on fly.io at
  `https://agent-bazaar-mcpgate.fly.dev` (placeholder — fly.io app not
  created yet). Browser → mcpgate needs the CORS support added in
  mcpgate (PR forthcoming); without it, cross-origin MCP calls are blocked.
- Demo LLM key: see `docs/demo-key.md`. No real key is in this repo, ever.

## Status

M1 skeleton (2026-11-16): networks + manifests + links + three-origin
bundles generate and smoke-test. Not yet verified in a real browser; the
bazaar UI (card wall, verification panel, trace, BYOK seeding) is M1/M2
work. MCP backends (mock data) are not written yet — the HOCON files
reference their planned mcpgate URLs.

## Disclosure (Devpost "New & Existing")

Agent Bazaar itself is newly created during the submission period
(Nov 2026–Jan 2027). It uses mcpgate, the author's pre-existing
MIT-licensed MCP gateway (v0.1.0), as its tool backend; the CORS support
added during the contest period is itself an enhancement shared back to
that project.

## License

MIT. See [LICENSE](LICENSE).
