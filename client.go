package main

import (
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
	idNum   = byte(0)
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	id    byte
	login string

	x     float64
	y     float64
	angle float64

	isUp    bool
	isDown  bool
	isRight bool
	isLeft  bool
	isShoot bool

	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		switch message[0] {
		case 0:
			c.isUp = true
		case 1:
			c.isDown = true
		case 2:
			c.isLeft = true
		case 3:
			c.isRight = true
		case 4:
			c.isShoot = true
		case 5:
			c.isUp = false
		case 6:
			c.isDown = false
		case 7:
			c.isLeft = false
		case 8:
			c.isRight = false
		case 9:
			c.isShoot = false
		}

		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) emulateClient() {
	var deltaTime = time.Now().UnixNano()
	var lastShotTime = time.Now().UnixNano()

	for {

		var deltaTime = time.Now().UnixNano() - deltaTime
		var delta = float64(deltaTime) / 1000000
		speed := float64(0)
		if c.isUp {
			speed += 2
		}
		if c.isDown {
			speed -= 2
		}
		sin := math.Sin(c.angle)
		cos := 1 - sin
		if sin < 0 {
			cos = 1 + sin
		}

		if c.isShoot && lastShotTime+500000 < time.Now().UnixNano() {

			bullet := Bullet{}
			if c.angle >= 0 && c.angle < 90 || c.angle < 270 && c.angle > 180 {
				bullet.x += c.x + delta*speed*cos
			} else {
				bullet.x = c.x - delta*speed*cos
			}
			bullet.y = c.y + delta*speed*sin

			Bullets[&bullet] = struct{}{}
		}

		if c.angle >= 0 && c.angle < 90 || c.angle < 270 && c.angle > 180 {
			c.x += delta * speed * cos
		} else {
			c.x -= delta * speed * cos
		}
		c.y += delta * speed * sin

		if c.isLeft {
			speed = delta * float64(90)
		}
		if c.isRight {
			speed = -delta * float64(90)
		}
		c.angle += speed
		if c.angle >= 360 {
			c.angle -= 360
		}

		<-time.NewTimer(time.Millisecond * 8).C
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub,
		conn:    conn,
		send:    make(chan []byte, 256),
		login:   r.FormValue("login"),
		x:       rand.Float64() * 32,
		y:       rand.Float64() * 32,
		isUp:    false,
		isDown:  false,
		isRight: false,
		isLeft:  false,
		isShoot: false,
	}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
	go client.emulateClient()
}
