package main

import (
	"log"

	"nhooyr.io/websocket"
)

// --- Hub with per-client outbound channels ---
// We keep a registry of active connections inside a dedicated goroutine.
// Each client owns a buffered outbound channel to decouple slow consumers.

type WebSocketClient struct {
	Conn *websocket.Conn
	Send chan []byte // outbound messages for this client
	Id   string      // optional identifier for logs
}

type WebSocketHub struct {
	Register   chan *WebSocketClient
	Unregister chan *WebSocketClient
	Broadcast  chan []byte
	clients    map[*WebSocketClient]struct{}
}

func NewWebSocketHub() *WebSocketHub {
	hub := &WebSocketHub{
		Register:   make(chan *WebSocketClient),
		Unregister: make(chan *WebSocketClient),
		Broadcast:  make(chan []byte, 1024), // global broadcast buffer
		clients:    make(map[*WebSocketClient]struct{}),
	}
	go hub.run()
	return hub
}

func (h *WebSocketHub) run() {
	for {
		select {
		case c := <-h.Register:
			h.clients[c] = struct{}{}
			log.Printf("ws: client registered (%s), total=%d", c.Id, len(h.clients))
		case c := <-h.Unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.Send) // signal writer to exit
				_ = c.Conn.Close(websocket.StatusNormalClosure, "bye")
				log.Printf("ws: client unregistered (%s), total=%d", c.Id, len(h.clients))
			}
		case msg := <-h.Broadcast:
			for c := range h.clients {
				// Non-blocking send; if buffer is full, drop the client to protect the hub
				select {
				case c.Send <- msg:
				default:
					log.Printf("ws: client buffer full, disconnecting (%s)", c.Id)
					delete(h.clients, c)
					close(c.Send)
					_ = c.Conn.Close(websocket.StatusPolicyViolation, "slow consumer")
				}
			}
		}
	}
}
