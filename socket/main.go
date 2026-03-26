package main

import (
	"log"
	"sync"

	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

// Hub con todos los navegadores conectados
type Hub struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

var hub = &Hub{clients: make(map[*websocket.Conn]bool)}

func (h *Hub) add(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[conn] = true
}

func (h *Hub) remove(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, conn)
}

func (h *Hub) broadcast(msg []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for conn := range h.clients {
		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			conn.Close()
			delete(h.clients, conn)
		}
	}
}

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		return true
	},
}

func wsHandler(ctx *fasthttp.RequestCtx) {
	upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		hub.add(conn)
		defer hub.remove(conn)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			log.Printf("Recibido de app: %s", msg)

			// 👇 Reenvía a todos los navegadores conectados
			hub.broadcast(msg)
		}
	})
}

func main() {
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/ws":
			wsHandler(ctx)
		default:
			ctx.SetStatusCode(fasthttp.StatusNotFound)
		}
	}

	log.Println("Servidor WS corriendo en :8081")
	if err := fasthttp.ListenAndServe(":8081", requestHandler); err != nil {
		log.Fatalf("Error: %s", err)
	}
}
