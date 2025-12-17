# Shared Canvas

This project is a minimal web server that:
- Provides web UI as SPA, embed into the binary
- Exposes a WebSocket and REST API

The UI provides a simple drawing canvas that can be shared with other users.
Every update on the local canvas is broadcast to all connected users, so they all see the same drawing.

## API

* `/` — returns the SPA HTML
* `/socket` — upgrades to WebSocket
* `/image` — returns a PNG image of the canvas

### WebSocket messages

#### init

Sent by the server to the SPA on WebSocket connected to initialise the canvas.

```json
{
  "method": "init",
  "params": {
    "w": 100,
    "h": 100,
    "p": "000111..."
  }
}
```

`w` and `h` are the canvas dimensions.

`p` is the pixel array encoded: `_` for transparent, `0` for black, `1` for white.

#### draw

Sent by SPA to the server and broadcast to all connected clients.

```json
{
  "method": "draw",
  "params": {
    "x": 10,
    "y": 20,
    "w": 3,
    "h": 3,
    "p": "___000___"
  }
}
```

`x`, `y` are the top-left corner of the rectangular area to draw.

`w` and `h` are the rectangular area dimensions.

`p` is the pixel array encoded: `_` for transparent, `0` for black, `1` for white.

#### user

Sent by the server to the SPA to notify the number of connected users.

```json
{
  "method": "user",
  "params": {
    "count": 10
  }
}
```

## Requirements

- Go 1.21+
- Node.js 22+

## Build

Using Makefile (recommended):

```bash
make build
```

This builds the Svelte SPA to `cmd/shared-canvas-server/web-dist/` and then compiles the Go server. The resulting binary is embedded with the SPA assets.

Manual steps:

```bash
# 1) Build the SPA (outputs to cmd/shared-canvas-server/web-dist)
cd webapp && npm ci && npm run build

# 2) Build the Go server
go build ./cmd/shared-canvas-server
```

## Run

```bash
./shared-canvas-server -port 8080 -width 576 -height 576 -image canvas.png
```

Then open your browser at:
```
http://localhost:8080/
```

### CLI options

- `-port` (int): Port to listen on (default: 8080)
- `-width` (int): Canvas width (default: 576)
- `-height` (int): Canvas height (default: 576)
- `-image` (string): Path to the image to load on startup and save to on shutdown (default: canvas.png)

## Implementation details

UI has two canvases. One above the other. 

The upper canvas is to draw the lines locally. 
It's transparent.
The line is drawn on the mouse move. 
Then the rectangular area of the line pixels is sent to the server.
The pixels are encoded, so only black, white and transparent pixels are sent.
All the grey-scaling appeared to smooth the line is lost.

The lower canvas is to display the shared image.
The same message with the line pixels is received from the server.
It's drawn on the lower canvas.
And the same rectangular area is cleared on the upper canvas to erase the local line.
