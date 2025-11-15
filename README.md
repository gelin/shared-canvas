# Shared Canvas — Go Web Server with Embedded Assets

This project is a minimal Golang web server that:
- Serves a small web UI (HTML/CSS/JS) compiled into the binary using `go:embed`.
- Exposes REST API endpoints.
- Allows configuring the listening port via CLI flag.
- Shuts down gracefully on SIGINT/SIGTERM.

## Requirements
- Go 1.21+

## Build
```bash
go build ./cmd/server
```
This produces a single binary (named `server` or `server.exe` depending on your OS) with the web assets embedded.

## Run
```bash
./server -port 8080
# or shorthand
./server -p 8080
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

## Embedded Web UI
The web UI is located under `cmd/server/web` and includes:
- `index.html`
- `style.css`
- `app.js`

These files are embedded into the binary via `//go:embed web/*` and served from memory. Unknown non-API routes fall back to `index.html` (SPA-friendly behavior).

## Project Layout
```
cmd/
  server/
    main.go       # server entrypoint, routes, embedding, graceful shutdown
    web/
      index.html
      style.css
      app.js
```

## Notes
- Logging middleware prints method, path, status, and duration for each request.
- The server attempts graceful shutdown with a 5s timeout upon receiving SIGINT/SIGTERM.
