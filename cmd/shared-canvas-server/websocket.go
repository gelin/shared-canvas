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
	defer func(c *websocket.Conn, code websocket.StatusCode, reason string) {
		err := c.Close(code, reason)
		if err != nil {
			log.Printf("ws close error: %v", err)
		}
	}(c, websocket.StatusNormalClosure, "bye")

	ctx := r.Context()

	// Send a greeting message
	//send := wsMessage{Method: "draw"}
	//_ = writeWSJSON(ctx, c, send)

	for {
		var msg wsMessage
		if err := readWSJSON(ctx, c, &msg); err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure || websocket.CloseStatus(err) == websocket.StatusGoingAway {
				return
			}
			log.Printf("ws read error: %v", err)
			return
		}
		log.Printf("ws received: %v", msg)
		// Echo/ack with the same message
		reply := msg
		if err := writeWSJSON(ctx, c, reply); err != nil {
			log.Printf("ws write error: %v", err)
			return
		}
		log.Printf("ws sent: %v", reply)
	}
}

func writeWSJSON(ctx context.Context, c *websocket.Conn, v any) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return c.Write(ctx, websocket.MessageText, mustJSON(v))
}

func mustJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
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
