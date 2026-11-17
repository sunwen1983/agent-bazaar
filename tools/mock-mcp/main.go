// Command mock-mcp is a minimal, deterministic MCP server that speaks
// JSON-RPC 2.0 over stdio. It exists so Agent Bazaar's three peer agents
// have real tools to call through mcpgate in the demo. One dataset per
// process, selected by argv[1]:
//
//	mock-mcp flights      # search_flights   (6 SFO->AUS round-trip options)
//	mock-mcp stays        # search_stays     (8 Austin hotels)
//	mock-mcp activities   # search_activities (10 Austin activities)
//
// The data is fixed sample inventory, not live availability: the same
// arguments always return the same answer, prices are realistic, and the
// cheapest combinations fit the $800 Austin-weekend budget in the demo
// scenario (2 adults, Fri-Sun). Every response says so in its "note".
//
// The wire protocol is modeled on mcpgate's examples/echo-backend: one
// JSON object per line on stdin/stdout, protocolVersion "2025-03-26".
// Logs go to stderr only; stdout carries protocol traffic.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type rpcRequest struct {
	JSONRPC string         `json:"jsonrpc"`
	ID      any            `json:"id,omitempty"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params,omitempty"`
}

type rpcResponse struct {
	JSONRPC string  `json:"jsonrpc"`
	ID      any     `json:"id,omitempty"`
	Result  any     `json:"result,omitempty"`
	Error   *rpcErr `json:"error,omitempty"`
}

type rpcErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// tool is one MCP tool: its schema plus the function that answers calls.
type tool struct {
	Name        string
	Description string
	Schema      map[string]any
	Call        func(args map[string]any) (any, *rpcErr)
}

type dataset struct {
	Name  string
	Tools []tool
}

func respond(w *bufio.Writer, id any, result any, err *rpcErr) {
	resp := rpcResponse{JSONRPC: "2.0", ID: id, Result: result, Error: err}
	line, _ := json.Marshal(resp)
	fmt.Fprintln(w, string(line))
	_ = w.Flush()
}

func handle(w *bufio.Writer, ds dataset, req rpcRequest) {
	switch req.Method {
	case "initialize":
		respond(w, req.ID, map[string]any{
			"protocolVersion": "2025-03-26",
			"capabilities":    map[string]any{"tools": map[string]any{}},
			"serverInfo":      map[string]any{"name": "agent-bazaar-mock-" + ds.Name, "version": "0.1.0"},
		}, nil)
	case "ping":
		respond(w, req.ID, map[string]any{}, nil)
	case "tools/list":
		tools := make([]any, 0, len(ds.Tools))
		for _, t := range ds.Tools {
			tools = append(tools, map[string]any{
				"name":        t.Name,
				"description": t.Description,
				"inputSchema": t.Schema,
			})
		}
		respond(w, req.ID, map[string]any{"tools": tools}, nil)
	case "tools/call":
		name, _ := req.Params["name"].(string)
		args, _ := req.Params["arguments"].(map[string]any)
		if args == nil {
			args = map[string]any{}
		}
		for _, t := range ds.Tools {
			if t.Name == name {
				result, err := t.Call(args)
				if err != nil {
					respond(w, req.ID, nil, err)
					return
				}
				text, _ := json.Marshal(result)
				respond(w, req.ID, map[string]any{
					"content":           []any{map[string]any{"type": "text", "text": string(text)}},
					"structuredContent": result,
				}, nil)
				return
			}
		}
		respond(w, req.ID, nil, &rpcErr{Code: -32602, Message: "unknown tool: " + name})
	case "notifications/initialized", "notifications/cancelled", "notifications/progress":
		// Notifications carry no id and get no response.
	default:
		respond(w, req.ID, nil, &rpcErr{Code: -32601, Message: "method not found: " + req.Method})
	}
}

// --- argument helpers (all forgiving: bad types fall back to defaults) ---

func strArg(args map[string]any, key, def string) string {
	if v, ok := args[key].(string); ok && v != "" {
		return v
	}
	return def
}

func floatArg(args map[string]any, key string, def float64) float64 {
	if v, ok := args[key].(float64); ok {
		return v
	}
	return def
}

func intArg(args map[string]any, key string, def int) int {
	if v, ok := args[key].(float64); ok {
		return int(v)
	}
	return def
}

func strSliceArg(args map[string]any, key string) []string {
	raw, ok := args[key].([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func containsFold(haystack, needle string) bool {
	return strings.Contains(strings.ToLower(haystack), strings.ToLower(needle))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: mock-mcp <flights|stays|activities>")
		os.Exit(2)
	}
	var ds dataset
	switch os.Args[1] {
	case "flights":
		ds = flightsDataset()
	case "stays":
		ds = staysDataset()
	case "activities":
		ds = activitiesDataset()
	default:
		fmt.Fprintf(os.Stderr, "unknown dataset %q (want flights|stays|activities)\n", os.Args[1])
		os.Exit(2)
	}
	fmt.Fprintf(os.Stderr, "mock-mcp %s: ready (deterministic sample data)\n", ds.Name)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	writer := bufio.NewWriter(os.Stdout)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var req rpcRequest
		if err := json.Unmarshal(line, &req); err != nil {
			respond(writer, nil, nil, &rpcErr{Code: -32700, Message: "parse error"})
			continue
		}
		handle(writer, ds, req)
	}
}
