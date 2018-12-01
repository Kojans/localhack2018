package main

import (
	"fmt"
	"time"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[byte]*Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[byte]*Client),
	}
}

func (h *Hub) run() {
	for {
		tik := time.NewTicker(8 * time.Millisecond)
		select {
		case client := <-h.register:

			//bs := make([]byte, 4)
			//binary.LittleEndian.PutUint32(bs, len(h.clients))
			//sendData := []byte{};
			for _, clientCurr := range h.clients {
				x := float64ToByte(clientCurr.x)
				y := float64ToByte(clientCurr.y)
				a := float64ToByte(clientCurr.angle)
				client.send <- []byte{clientCurr.id, 0,
					x[0], x[1], x[2], x[3], x[4], x[5], x[6], x[7],
					y[0], y[1], y[2], y[3], y[4], y[5], y[6], y[7],
					a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7]}
			}
			h.clients[idNum] = client
			idNum++
		case client := <-h.unregister:
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients, client.id)
				close(client.send)
			}
		case message := <-h.broadcast:
			for _, client := range h.clients {
				fmt.Println(222)
				select {
				case client.send <- message:
				default:
					fmt.Println(111)
					close(client.send)
					delete(h.clients, client.id)
				}
			}
		case <-tik.C:
			for _, client := range h.clients {
				for bullet := range Bullets {
					if bullet.x > client.x-0.5 && bullet.x < client.x+0.5 &&
						bullet.y > client.y-0.5 && bullet.y < client.y+0.5 {
						h.unregister <- client

						delete(Bullets, bullet)

						x := float64ToByte(client.x)
						y := float64ToByte(client.y)
						a := float64ToByte(client.angle)
						h.broadcast <- []byte{client.id, 1,
							x[0], x[1], x[2], x[3], x[4], x[5], x[6], x[7],
							y[0], y[1], y[2], y[3], y[4], y[5], y[6], y[7],
							a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7]}
					}
				}
			}
		}
	}
}
