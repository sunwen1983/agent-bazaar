# syntax=docker/dockerfile:1
# Agent Bazaar mcpgate image.
#
# Multi-stage: builds mcpgate (pinned commit 9bb4a88, includes the CORS
# middleware from PR #20) and the mock-mcp backends from this repo, then
# ships both binaries plus the gateway config template. The demo API key
# is NOT baked in: entrypoint.sh substitutes it from the DEMO_API_KEY
# environment variable (a fly secret) at container start.
#
# Build context note: this Dockerfile clones both repos from GitHub, so
# it is self-contained — `cd mcpgate && fly deploy` works without the
# rest of the agent-bazaar checkout.
ARG MCPGATE_REF=9bb4a88
ARG AGENT_BAZAAR_REF=main

FROM golang:1.27-bookworm AS build
ARG MCPGATE_REF
ARG AGENT_BAZAAR_REF
RUN git clone https://github.com/sunwen1983/mcpgate /src/mcpgate \
 && cd /src/mcpgate && git checkout ${MCPGATE_REF} \
 && go build -o /out/mcpgate ./cmd/mcpgate
RUN git clone --depth 1 --branch ${AGENT_BAZAAR_REF} https://github.com/sunwen1983/agent-bazaar /src/agent-bazaar \
 && cd /src/agent-bazaar/tools/mock-mcp && go build -o /out/mock-mcp .

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
 && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=build /out/mcpgate /out/mock-mcp /app/
COPY config.yaml /app/config.yaml.template
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh
EXPOSE 8080
ENTRYPOINT ["/app/entrypoint.sh"]
