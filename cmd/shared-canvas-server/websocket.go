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

var hub = NewWebSocketHub()

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Optional origin check can be added here
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // for the demo; consider proper origin checks in production
	})
	if err != nil {
		log.Printf("ws accept error: %v", err)
		return
	}
	// Prepare the client and register to hub
	client := &WebSocketClient{
		Conn: c,
		Send: make(chan []byte, 256), // per-client buffer; tune as needed
		Id:   r.RemoteAddr,
	}
	hub.Register <- client

	// Send init message
	err = writeWSJSON(r.Context(), c, imgHolder.GetImageAsInitMessage())
	if err != nil {
		log.Printf("ws init error (%s): %v", client.Id, err)
	}

	// Start a writer goroutine for this client
	go writePump(r.Context(), client)

	// Reader loop: read messages and broadcast
	for {
		var msg *DrawMessage
		if err := readWSJSON(r.Context(), c, &msg); err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure || websocket.CloseStatus(err) == websocket.StatusGoingAway {
				break
			}
			log.Printf("ws read error (%s): %v", client.Id, err)
			break
		}
		log.Printf("ws received (%s): %v", client.Id, msg)
		broadcastWSJSON(r.Context(), msg)
	}

	// ensure unregistration
	hub.Unregister <- client
}

// writePump consumes the client's send channel and writes frames to the socket.
// Exits when the channel is closed or a write error occurs.
func writePump(ctx context.Context, c *WebSocketClient) {
	for msg := range c.Send {
		wctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		err := c.Conn.Write(wctx, websocket.MessageText, msg)
		cancel()
		if err != nil {
			log.Printf("ws write error (%s): %v", c.Id, err)
			break
		}
	}
	// Ensure connection is closed on writer exit; hub.Unregister will also try closing, but double-close is fine
	_ = c.Conn.Close(websocket.StatusNormalClosure, "writer exit")
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
	drawMsg, ok := msg.(*DrawMessage)
	if ok {
		imgHolder.Draw <- drawMsg
	} else {
		log.Fatalf("unexpected message type: %T\n", msg)
	}
	data := mustJSON(msg)
	select {
	case hub.Broadcast <- data:
	default:
		// if a global broadcast is saturated, attempt a bounded wait; otherwise drop
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		select {
		case hub.Broadcast <- data:
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

func writeWSJSON(ctx context.Context, c *websocket.Conn, v any) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return c.Write(ctx, websocket.MessageText, mustJSON(v))
}
