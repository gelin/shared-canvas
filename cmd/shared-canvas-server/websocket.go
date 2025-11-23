package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"nhooyr.io/websocket"
)

type wsMessage struct {
	Method string      `json:"method"`
	Params interface{} `json:"params,omitempty"`
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Optional origin check can be added here
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // for the demo; consider proper origin checks in production
	})
	if err != nil {
		log.Printf("ws accept error: %v", err)
		return
	}
	// Prepare the client and register to hub
	client := &wsClient{
		conn: c,
		send: make(chan []byte, 256), // per-client buffer; tune as needed
		id:   r.RemoteAddr,
	}
	hub.register <- client

	// Start writer goroutine for this client
	go writePump(r.Context(), client)

	// Reader loop: read messages and broadcast
	for {
		var msg wsMessage
		if err := readWSJSON(r.Context(), c, &msg); err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure || websocket.CloseStatus(err) == websocket.StatusGoingAway {
				break
			}
			log.Printf("ws read error (%s): %v", client.id, err)
			break
		}
		log.Printf("ws received (%s): %v", client.id, msg)
		broadcastWSJSON(r.Context(), msg)
	}

	// ensure unregistration
	hub.unregister <- client
}

// writePump consumes the client's send channel and writes frames to the socket.
// Exits when the channel is closed or a write error occurs.
func writePump(ctx context.Context, c *wsClient) {
	for msg := range c.send {
		wctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		err := c.conn.Write(wctx, websocket.MessageText, msg)
		cancel()
		if err != nil {
			log.Printf("ws write error (%s): %v", c.id, err)
			break
		}
	}
	// Ensure connection is closed on writer exit; hub.unregister will also try closing, but double-close is fine
	_ = c.conn.Close(websocket.StatusNormalClosure, "writer exit")
}

func readWSJSON(ctx context.Context, c *websocket.Conn, v any) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	typ, data, err := c.Read(ctx)
	if err != nil {
		return err
	}
	if typ != websocket.MessageText {
		return fmt.Errorf("unexpected message type: %v", typ)
	}
	return json.Unmarshal(data, v)
}

func broadcastWSJSON(ctx context.Context, msg any) {
	// marshal once and fan out
	data := mustJSON(msg)
	select {
	case hub.broadcast <- data:
	default:
		// if a global broadcast is saturated, attempt a bounded wait; otherwise drop
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		select {
		case hub.broadcast <- data:
			cancel()
		case <-ctx.Done():
			cancel()
			log.Printf("ws: global broadcast channel saturated, dropping message")
		}
	}
}

func mustJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}
