# Shared Canvas — Go Web Server with Embedded SvelteKit (Svelte 5) SPA

This project is a minimal Golang web server that:
- Serves a SvelteKit (Svelte 5) SPA compiled into the binary using `go:embed` (no runtime Node.js needed to run the binary)
- Exposes REST API endpoints
- Allows configuring the listening port via CLI flag
- Shuts down gracefully on SIGINT/SIGTERM

## Requirements
- Go 1.21+
- Node.js 18+ (only required for building the frontend)

## Build (frontend + backend)
Using Makefile (recommended):
```bash
make build
```
This builds the Svelte SPA to `cmd/shared-canvas-server/web-dist/` and then compiles the Go server. The resulting binary is embedded with the SPA assets.

Manual steps:
```bash
# 1) Build the SPA (outputs to cmd/shared-canvas-server/web-dist)
(cd webapp && npm ci && npm run build)

# 2) Build the Go server
go build ./cmd/shared-canvas-server
```

## Run
```bash
./shared-canvas-server -port 8080
```
Then open your browser at:
```
http://localhost:8080/
```

CLI options:
- `-port` (int): Port to listen on (default: 8080)
- `-p` (int): Shorthand for `-port` (default: 8080)

## REST API
- `GET /api/health` → `{ "status": "ok", "service": "shared-canvas", "time": "RFC3339Nano" }`
- `GET /api/time` → `{ "now": "RFC3339Nano", "epoch": 173... }`

## WebSocket
- Endpoint: `GET /ws` (upgrades to WebSocket)
- Protocol: JSON messages `{ type: string, payload?: any, time?: string }`
- Behavior: server sends an initial `{ type: "welcome", payload: { message: "connected" }, time }` then echoes back any JSON you send wrapped in `{ type: "ack", payload: <your message>, time }`.

Dev usage:
- Terminal A: `go run ./cmd/shared-canvas-server -port 8080`
- Terminal B: `cd webapp && npm run dev`
- Open http://localhost:5173/ws and click Connect. The dev server proxies `ws://localhost:5173/ws` → `http://localhost:8080/ws`.

Production usage:
- `make build` then run the binary (e.g., `./shared-canvas-server -p 8080`)
- Open http://localhost:8080/ws and connect.

Notes:
- Go dependency used: `nhooyr.io/websocket`. If building manually the first time, run `go mod tidy` to fetch it.

## Frontend (Svelte 5 SPA)
- Source: `webapp/` (Vite + Svelte)
- Dev server:
  ```bash
  # Terminal A: run the Go API
  go run ./cmd/shared-canvas-server -port 8080
  # Terminal B: run SvelteKit dev server (with API proxy)
  cd webapp && npm run dev
  ```
  SvelteKit dev server runs on http://localhost:5173 and proxies `/api/*` to http://localhost:8080.
- Production build artifacts are written to `cmd/shared-canvas-server/web-dist/` and embedded into the Go binary.
- Unknown non-API routes fall back to `index.html` (SPA routing).

## Project Layout
```
cmd/
  shared-canvas-server/
    main.go        # server entrypoint, routes, embedding, graceful shutdown
    web-dist/      # built SPA assets (embedded)
webapp/            # Svelte 5 app (Vite)
  src/
    App.svelte
    routes/
      Home.svelte
      Health.svelte
      Time.svelte
  vite.config.ts
  package.json
```

## Notes
- Logging middleware prints method, path, status, and duration for each request.
- The server attempts graceful shutdown with a 5s timeout upon receiving SIGINT/SIGTERM.
