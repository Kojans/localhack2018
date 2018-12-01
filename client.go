package main

import (
	"websocket-master"
	"log"
	"time"
	"math"
	"encoding/binary"
)

func (c *Client) reader() {
	c.conn.SetReadLimit(512)
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		//log.Println(message)
		//W - 87, D - 68, A - 65, S - 83
		switch message[0] {
		case 1:
			switch message[1] {
			case 87:
				if c.goV > 0 {
					continue
				}
				c.goV++
				break
			case 68:
				if c.goH > 0 {
					continue
				}
				c.goH++
				break
			case 65:
				if c.goH < 0 {
					continue
				}
				c.goH--
				break
			case 83:
				if c.goV < 0 {
					continue
				}
				c.goV--
				break
			}
			break
		case 2:
			switch message[1] {
			case 83:
				if c.goV > 0 {
					continue
				}
				c.goV++
				break
			case 65:
				if c.goH > 0 {
					continue
				}
				c.goH++
				break
			case 68:
				if c.goH < 0 {
					continue
				}
				c.goH--
				break
			case 87:
				if c.goV < 0 {
					continue
				}
				c.goV--
				break
			}
			break
		case 3:
			/*if len(message) != 5 {
				return
			}
			for i := 0; i < 4; i++ {
				c.angleB[i] = message[i+1]
			}
			c.angle = math.Float32frombits(binary.LittleEndian.Uint32(message[1:]))*/
			break
		case 4:
			if c.shoting_now {
				break
			}
			c.shoting_now = true
			c.shoting <- true
			break
		case 5:
			if !c.shoting_now {
				break
			}
			c.shoting_now = false
			c.shoting <- false
			break
		case 6:
			if len(message) != 5 {
				return
			}
			for i := 0; i < 4; i++ {
				c.a_gun_angleB[i] = message[i+1]
			}
			c.a_gun_angle = math.Float32frombits(binary.LittleEndian.Uint32(message[1:]))
			break
		case 7:
			//log.Println(c.base)
			if c.base {
				if c.x <= 1000 && c.y <= 1000 {
					c.room.Base_A.Points += uint32(c.res[0] + c.res[1]*7 + c.res[2]*35)
					c.points += uint32(c.res[0] + c.res[1]*7 + c.res[2]*35)
					c.res = [3]uint16{0, 0, 0}
				}
			} else {
				if c.x >= 15000 && c.y >= 15000 {
					c.room.Base_B.Points += uint32(c.res[0] + c.res[1]*7 + c.res[2]*35)
					c.points += uint32(c.res[0] + c.res[1]*7 + c.res[2]*35)
					c.res = [3]uint16{0, 0, 0}
				}
			}
			forSend := make([]byte, 7)
			forSend[0] = 6

			c.send <- forSend[:]
			break
		case 8:
			forSend := make([]byte, 1)

			forSend[0] = 11

			MG := make([]byte, c.Ship.Main_Guns*2+1)
			MG[0] = c.Ship.Main_Guns
			i := 1
			for q, s := range c.Main_Guns {
				MG[i] = q
				MG[i+1] = s.t
				i += 2
			}
			forSend = append(forSend, MG...)

			MG = make([]byte, c.Ship.Additional_Guns*2+1)
			MG[0] = c.Ship.Additional_Guns
			i = 1
			for q, s := range c.Additional_Guns {
				MG[i] = q
				MG[i+1] = s.t
				i += 2
			}
			forSend = append(forSend, MG...)

			MG = make([]byte, c.Ship.Accelerators*2+1)
			MG[0] = c.Ship.Accelerators
			i = 1
			for q, s := range c.Accelerators {
				MG[i] = q
				MG[i+1] = s.t
				i += 2
			}
			forSend = append(forSend, MG...)

			MG = make([]byte, c.Ship.Engine*2+1)
			MG[0] = c.Ship.Engine
			i = 1
			for q, s := range c.Engine {
				MG[i] = q
				MG[i+1] = s.t
				i += 2
			}
			forSend = append(forSend, MG...)

			MG = make([]byte, c.Ship.Generators*2+1)
			MG[0] = c.Ship.Generators
			i = 1
			for q, s := range c.Generators {
				MG[i] = q
				MG[i+1] = s.t
				i += 2
			}
			forSend = append(forSend, MG...)

			c.send <- forSend[:]
			break
		case 9:
			if uint8(message[1]) < 3 {
				if message[1] == c.Ship.Type {
					continue
				}
				if message[1] == 0 {
					c.Ship = Ships[0]
					for i, _ := range c.Main_Guns {
						delete(c.Main_Guns, i)
					}
					c.Main_Guns[0] = &GunOnShip{80,160, 0}
					c.SetShip();
					continue
				}
				if message[1] == 1 {
					if c.points >= 1 {
						c.Ship = Ships[1]
						for i, _ := range c.Main_Guns {
							delete(c.Main_Guns, i)
						}
						c.Main_Guns[0] = &GunOnShip{80,160, 0}
						c.points -= 1
						c.SetShip();
					}
					continue
				}
				if message[1] == 2 {
					if c.points >= 1 {
						c.Ship = Ships[2]
						for i, _ := range c.Main_Guns {
							delete(c.Main_Guns, i)
						}
						c.Main_Guns[0] = &GunOnShip{20,240, 0}
						c.Main_Guns[1] = &GunOnShip{110,240, 0}
						c.points -= 1
						c.SetShip();
					}
					continue
				}
			}
			break
		}
	}
}

func (c *Client) writer() {
	for {
		message, ok := <-c.send
		c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		if !ok {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		w, err := c.conn.NextWriter(websocket.BinaryMessage)
		if err != nil {
			return
		}
		w.Write(message)
		if err := w.Close(); err != nil {
			return
		}
	}
}

type Shell struct {
	x, y   float64
	close  bool
	damage uint16
	t      bool
}

func NewClient(ws *websocket.Conn, name string, x, y float64, room *Room, base bool) *Client {
	return &Client{
		ws,
		name,
		make(chan []byte, 32),
		x,
		y,
		0,
		0,
		0,
		0,
		[]byte{0, 0, 0, 0},
		0,
		[]byte{0, 0, 0, 0},
		10,
		150,
		150,
		10,
		10,
		make(chan bool, 2),
		false,
		50,
		false,
		10,
		50,
		[3]uint16{0, 0, 0},
		room,
		base,
		Ships[0],
		make(map[uint8]*GunOnShip),
		make(map[uint8]*GunOnShip),
		make(map[uint8]*Generator),
		make(map[uint8]*Accelerator),
		make(map[uint8]*Engine),
	}
}

func (c *Client) SetHP() {
	hp := uint16(0)
	hp += c.Ship.HP
	c.max_hp = hp
}

func (c *Client) SetShield() {
	hp := uint16(0)
	hp += uint16(c.Ship.Shield)
	c.max_shield = hp
}

func (c *Client) SetShip() {
	c.SetSpeed()
	c.SetHealth()
	c.res = [3]uint16{0, 0, 0}
	c.space = c.Ship.Space
}

func (c *Client) SetHealth() {
	c.max_hp = c.Ship.HP
	c.hp = c.Ship.HP
}
func (c *Client) SetSpeed() {
	s := float64(0)
	for _, o := range c.Engine {
		s += o.Speed
	}
	c.speed = s * c.Ship.Speed
}

type Client struct {
	conn *websocket.Conn
	name string
	send chan []byte
	x, y float64
	goH  int8
	goV  int8

	points uint32

	a_gun_angle  float32
	a_gun_angleB []byte

	angle  float32
	angleB []byte
	speed  float64

	hp     uint16
	max_hp uint16

	shield     uint16
	max_shield uint16

	shoting      chan bool
	shoting_now  bool
	shell_speed  float64
	close        bool
	shell_damage uint16

	space uint16
	res   [3]uint16

	room *Room
	base bool
	Ship *Ship

	Main_Guns       map[uint8]*GunOnShip
	Additional_Guns map[uint8]*GunOnShip
	Generators      map[uint8]*Generator
	Accelerators    map[uint8]*Accelerator
	Engine          map[uint8]*Engine
}
