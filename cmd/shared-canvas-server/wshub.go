package main

import (
	"log"

	"nhooyr.io/websocket"
)

// --- Hub with per-client outbound channels ---
// We keep a registry of active connections inside a dedicated goroutine.
// Each client owns a buffered outbound channel to decouple slow consumers.

type wsClient struct {
	conn *websocket.Conn
	send chan []byte // outbound messages for this client
	id   string      // optional identifier for logs
}

type wsHub struct {
	register   chan *wsClient
	unregister chan *wsClient
	broadcast  chan []byte
	clients    map[*wsClient]struct{}
}

var hub *wsHub

func initWsHub() {
	hub = &wsHub{
		register:   make(chan *wsClient),
		unregister: make(chan *wsClient),
		broadcast:  make(chan []byte, 1024), // global broadcast buffer
		clients:    make(map[*wsClient]struct{}),
	}
	go hub.run()
}

func (h *wsHub) run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = struct{}{}
			log.Printf("ws: client registered (%s), total=%d", c.id, len(h.clients))
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send) // signal writer to exit
				_ = c.conn.Close(websocket.StatusNormalClosure, "bye")
				log.Printf("ws: client unregistered (%s), total=%d", c.id, len(h.clients))
			}
		case msg := <-h.broadcast:
			for c := range h.clients {
				// Non-blocking send; if buffer is full, drop the client to protect the hub
				select {
				case c.send <- msg:
				default:
					log.Printf("ws: client buffer full, disconnecting (%s)", c.id)
					delete(h.clients, c)
					close(c.send)
					_ = c.conn.Close(websocket.StatusPolicyViolation, "slow consumer")
				}
			}
		}
	}
}
