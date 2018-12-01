package main

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
		select {
		case client := <-h.register:
			for _, clientCurr := range h.clients {
				x := float64ToByte(clientCurr.x)
				y := float64ToByte(clientCurr.y)
				a := float64ToByte(clientCurr.angle)
				client.send <- append([]byte{}, clientCurr.id,
					x[0], x[1], x[2], x[3], x[4], x[5], x[6], x[7],
					y[0], y[1], y[2], y[3], y[4], y[5], y[6], y[7],
					a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7])
			}
			h.clients[idNum] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.id]; ok {
				delete(h.clients, client.id)
				close(client.send)
			}
		case message := <-h.broadcast:
			for _, client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client.id)
				}
			}
		}
	}
}
