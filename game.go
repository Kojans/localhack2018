package main

import (
	"time"
	"math"
	"encoding/binary"
	"log"
)

func (c *Client) monitoring() {
	CurrRoom.RW.RLock()

	CurrRoom.RW.RUnlock()
}

func (c *Client) Shooting() {
	t := time.NewTicker(time.Second * 2)
	for {
		select {

		case m := <-c.shoting:
			if m {
				func(c *Client, t *time.Ticker) {
					for {
						if c.close {
							close(c.shoting)
							return
						}
						select {
						case <-c.shoting:
							return
							break
						case <-t.C:
							for _, o := range c.Main_Guns {
								c.room.ShellMut.Lock()
								a := c.room.ShellId
								c.room.ShellId += 1
								c.room.ShellMut.Unlock()

								AngleCos := -1 * math.Cos(float64(c.angle))
								AngleSin := 1 * math.Sin(float64(c.angle))
								x1 := float64(c.x) + AngleSin*(float64(o.y)-80) + math.Sin(float64(c.angle)+math.Pi/2)*(float64(o.x)-65)
								if x1 < 0 || x1 > 16000 {
									continue
								}
								y1 := float64(c.y) + AngleCos*(float64(o.y)-80) - math.Cos(float64(c.angle)+math.Pi/2)*(float64(o.x)-65)
								if y1 < 0 || y1 > 16000 {
									continue
								}
								s := &Shell{x1, y1, false, uint16(Guns[o.t].Damage), c.base}

								//log.Println(s.x, s.y)

								c.room.ShellsMut.Lock()
								c.room.Shells[a] = s
								c.room.ShellsMut.Unlock()

								var nowX float64
								var nowY float64

								if AngleCos > 0 {
									nowY = AngleCos * c.shell_speed
								}
								if AngleCos < 0 {
									nowY = -1 * AngleCos * c.shell_speed
								}
								if AngleSin > 0 {
									nowX = AngleSin * c.shell_speed
								}
								if AngleSin < 0 {
									nowX = -1 * AngleSin * c.shell_speed
								}

								go Shoot(c.shell_speed, c.room, s, nowX, nowY, a, AngleSin, AngleCos)
							}
							break
						}
					}
				}(c, t)
			}
			break
		default:
			if c.close {
				close(c.shoting)
				return
			}
			break
		}
	}
}

func Shoot(speed float64, room *Room, s *Shell, nowX float64, nowY float64, a uint32, AngleSin, AngleCos float64) {

	t := time.NewTimer(time.Second * 3)
	f := time.NewTicker(time.Second / 120)
	for {
		if s.close {
			room.ShellsMut.Lock()
			delete(room.Shells, a)
			room.ShellsMut.Unlock()
			forSend := make([]byte, 5)
			forSend[0] = 3
			sd := make([]byte, 4)
			binary.LittleEndian.PutUint32(sd, a)
			forSend[1] = sd[3]
			forSend[2] = sd[2]
			forSend[3] = sd[1]
			forSend[4] = sd[0]
			CurrRoom.RW.RLock()
			for _, u := range CurrRoom.Clients {
				u.send <- forSend
			}
			CurrRoom.RW.RUnlock()
			return
		}
		select {
		case <-t.C:
			room.ShellsMut.Lock()
			delete(room.Shells, a)
			room.ShellsMut.Unlock()
			forSend := make([]byte, 5)
			forSend[0] = 3
			sd := make([]byte, 4)
			binary.LittleEndian.PutUint32(sd, a)
			forSend[1] = sd[3]
			forSend[2] = sd[2]
			forSend[3] = sd[1]
			forSend[4] = sd[0]
			CurrRoom.RW.RLock()
			for _, u := range CurrRoom.Clients {
				u.send <- forSend
			}
			CurrRoom.RW.RUnlock()
			return
		case <-f.C:
			if AngleCos > 0 {
				if nowY > room.Size_Y-s.y {
					room.ShellsMut.Lock()
					delete(room.Shells, a)
					room.ShellsMut.Unlock()
					forSend := make([]byte, 5)
					forSend[0] = 3
					sd := make([]byte, 4)
					binary.LittleEndian.PutUint32(sd, a)
					forSend[1] = sd[3]
					forSend[2] = sd[2]
					forSend[3] = sd[1]
					forSend[4] = sd[0]
					CurrRoom.RW.RLock()
					for _, u := range CurrRoom.Clients {
						u.send <- forSend
					}
					CurrRoom.RW.RUnlock()
					return
				} else {
					s.y += nowY
				}
			}
			if AngleCos < 0 {
				if nowY > s.y {
					room.ShellsMut.Lock()
					delete(room.Shells, a)
					room.ShellsMut.Unlock()
					forSend := make([]byte, 5)
					forSend[0] = 3
					sd := make([]byte, 4)
					binary.LittleEndian.PutUint32(sd, a)
					forSend[1] = sd[3]
					forSend[2] = sd[2]
					forSend[3] = sd[1]
					forSend[4] = sd[0]
					CurrRoom.RW.RLock()
					for _, u := range CurrRoom.Clients {
						u.send <- forSend
					}
					CurrRoom.RW.RUnlock()
					return
				} else {
					s.y -= nowY
				}
			}
			if AngleSin > 0 {
				if nowX > room.Size_X-s.x {
					room.ShellsMut.Lock()
					delete(room.Shells, a)
					room.ShellsMut.Unlock()
					forSend := make([]byte, 5)
					forSend[0] = 3
					sd := make([]byte, 4)
					binary.LittleEndian.PutUint32(sd, a)
					forSend[1] = sd[3]
					forSend[2] = sd[2]
					forSend[3] = sd[1]
					forSend[4] = sd[0]
					CurrRoom.RW.RLock()
					for _, u := range CurrRoom.Clients {
						u.send <- forSend
					}
					CurrRoom.RW.RUnlock()
					return
				} else {
					s.x += nowX
				}
			}
			if AngleSin < 0 {
				if nowX > s.x {
					room.ShellsMut.Lock()
					delete(room.Shells, a)
					room.ShellsMut.Unlock()
					forSend := make([]byte, 5)
					forSend[0] = 3
					sd := make([]byte, 4)
					binary.LittleEndian.PutUint32(sd, a)
					forSend[1] = sd[3]
					forSend[2] = sd[2]
					forSend[3] = sd[1]
					forSend[4] = sd[0]
					CurrRoom.RW.RLock()
					for _, u := range CurrRoom.Clients {
						u.send <- forSend
					}
					CurrRoom.RW.RUnlock()
					return
				} else {
					s.x -= nowX
				}
			}
			break
		}
	}
}

func GetByteShell(c *Shell, q uint32) []byte {
	forSend := make([]byte, 21)
	forSend[0] = 2
	a := make([]byte, 4)
	binary.LittleEndian.PutUint32(a, q)
	forSend[1] = a[3]
	forSend[2] = a[2]
	forSend[3] = a[1]
	forSend[4] = a[0]
	f := Float64bytes(c.x)
	forSend[5] = f[0]
	forSend[6] = f[1]
	forSend[7] = f[2]
	forSend[8] = f[3]
	forSend[9] = f[4]
	forSend[10] = f[5]
	forSend[11] = f[6]
	forSend[12] = f[7]
	f = Float64bytes(c.y)
	forSend[13] = f[0]
	forSend[14] = f[1]
	forSend[15] = f[2]
	forSend[16] = f[3]
	forSend[17] = f[4]
	forSend[18] = f[5]
	forSend[19] = f[6]
	forSend[20] = f[7]
	return forSend[:]
}

func GetByteEat(c *Eat, q uint32) []byte {
	forSend := make([]byte, 21)
	forSend[0] = 5
	a := make([]byte, 4)
	binary.LittleEndian.PutUint32(a, q)
	forSend[1] = a[3]
	forSend[2] = a[2]
	forSend[3] = a[1]
	forSend[4] = a[0]
	f := Float64bytes(c.x)
	forSend[5] = f[0]
	forSend[6] = f[1]
	forSend[7] = f[2]
	forSend[8] = f[3]
	forSend[9] = f[4]
	forSend[10] = f[5]
	forSend[11] = f[6]
	forSend[12] = f[7]
	f = Float64bytes(c.y)
	forSend[13] = f[0]
	forSend[14] = f[1]
	forSend[15] = f[2]
	forSend[16] = f[3]
	forSend[17] = f[4]
	forSend[18] = f[5]
	forSend[19] = f[6]
	forSend[20] = f[7]
	return forSend[:]
}

func (c *Client) RecoveryShield() {

	t := time.NewTicker(time.Second * 3)
	for {
		<-t.C
		if c.shield == c.max_shield {
			continue
		}
		if uint16(c.Ship.Shield_recovery) < c.max_shield-c.shield {
			c.shield += uint16(c.Ship.Shield_recovery)
		} else {
			c.shield = c.max_shield
		}
	}
}

func (c *Client) collision() {
	t := time.NewTicker(time.Second / 150)
	for {
		<-t.C
		/*if c.base {
			if c.x < 1000 && c.y < 1000 {
				c.room.Base_A.Points += c.points
				c.points = 0
			}
		} else {
			if c.x > 15000 && c.y > 15000 {
				c.room.Base_A.Points += c.points
				c.points = 0
			}
		}*/
		c.room.EatMut.RLock()
		for o, u := range c.room.Eats {
			distX := math.Abs(float64(u.x) - float64(c.x))
			distY := math.Abs(float64(u.y) - float64(c.y))
			if distX > 90 || distY > 90 || c.space <= c.res[0]+c.res[1]+c.res[2] {
				continue
			}
			c.res[u.t]++
			forSend := make([]byte, 7)
			forSend[0] = 6
			a := make([]byte, 2)
			binary.LittleEndian.PutUint16(a, c.res[0])
			forSend[1] = a[1]
			forSend[2] = a[0]
			binary.LittleEndian.PutUint16(a, c.res[1])
			forSend[3] = a[1]
			forSend[4] = a[0]
			binary.LittleEndian.PutUint16(a, c.res[2])
			forSend[5] = a[1]
			forSend[6] = a[0]

			c.send <- forSend[:]

			c.room.EatMut.RUnlock()
			c.room.EatMut.Lock()
			delete(c.room.Eats, o)
			c.room.EatMut.Unlock()
			c.room.EatMut.RLock()

			a = make([]byte, 4)
			forSend = make([]byte, 5)
			forSend[0] = 10
			binary.LittleEndian.PutUint32(a, o)
			forSend[1] = a[3]
			forSend[2] = a[2]
			forSend[3] = a[1]
			forSend[4] = a[0]

			log.Println(1)
			CurrRoom.RW.RLock()
			for _, u := range CurrRoom.Clients {
				u.send <- forSend
			}
			CurrRoom.RW.RUnlock()
			log.Println(2)
		}
		c.room.EatMut.RUnlock()
		c.room.ShellsMut.RLock()
		for _, u := range c.room.Shells {
			if u.t == c.base {
				continue
			}
			distX := math.Abs(float64(u.x) - float64(c.x))
			distY := math.Abs(float64(u.y) - float64(c.y))
			if distX > 110 || distY > 110 {
				continue
			}
			if c.shield > u.damage {
				c.shield -= u.damage
				u.damage = 0
			} else {
				u.damage -= c.shield
				c.shield = 0;
			}
			if c.hp > u.damage {
				c.hp -= u.damage
			} else {
				l := len(c.name)
				forSend := make([]byte, l+2)
				forSend[0] = 4
				forSend[1] = byte(l)
				for i := 2; i < l+2; i++ {
					forSend[i] = c.name[i-2]
				}
				CurrRoom.RW.Lock()
				for _, u := range CurrRoom.Clients {
					u.send <- forSend
				}
				delete(CurrRoom.Clients, c.name)
				CurrRoom.RW.Unlock()

				c.close = true
				c.room.ShellsMut.RUnlock()
				u.close = true

				return
			}

			u.close = true
		}
		c.room.ShellsMut.RUnlock()
	}
}
func (c *Client) start() {
	t1 := time.NewTicker(time.Second / 120)
	t2 := time.NewTicker(time.Second / 60)
	log.Println("Some connect")
	for {
		<-t1.C
		if c.close {
			if c.base {
				c.room.AMut.Lock()
				delete(c.room.Command_A, c)
				c.room.AMut.Unlock()
			} else {
				c.room.BMut.Lock()
				delete(c.room.Command_B, c)
				c.room.BMut.Unlock()
			}
			c.conn.Close()
			close(c.send)
			return
		}
		if c.goH > 0 {
			c.angle += math.Pi / 120
			if c.angle > math.Pi*2 {
				c.angle -= math.Pi * 2
			}
			c.angleB = Float32bytes(c.angle)
		}
		if c.goH < 0 {
			c.angle -= math.Pi / 120
			if c.angle < -math.Pi*2 {
				c.angle += math.Pi * 2
			}
			c.angleB = Float32bytes(c.angle)
		}
		if c.goV > 0 {
			Angle := -1 * c.speed * math.Cos(float64(c.angle))
			if Angle > 0 {
				nowY := Angle
				if nowY > c.room.Size_Y-c.y {
					c.y = c.room.Size_Y
				} else {
					c.y += nowY
				}
			}
			if Angle < 0 {
				nowY := -1 * Angle
				if nowY > c.y {
					c.y = 0
				} else {
					c.y -= nowY
				}
			}
			Angle = c.speed * math.Sin(float64(c.angle))
			if Angle > 0 {
				nowX :=Angle
				if nowX > c.room.Size_X-c.x {
					c.x = c.room.Size_X
				} else {
					c.x += nowX
				}
			}
			if Angle < 0 {
				nowX := -1 * Angle
				if nowX > c.x {
					c.x = 0
				} else {
					c.x -= nowX
				}
			}
		}
		if c.goV < 0 {
			Angle := -1 * c.speed / 2 * math.Cos(float64(c.angle))
			if Angle > 0 {
				nowY := Angle
				if nowY > c.y {
					c.y = 0
				} else {
					c.y -= nowY
				}
			}
			if Angle < 0 {
				nowY := -1 * Angle
				if nowY > c.room.Size_Y-c.y {
					c.y = c.room.Size_Y
				} else {
					c.y += nowY
				}
			}
			Angle = c.speed / 2 * math.Sin(float64(c.angle))
			if Angle > 0 {
				nowX := Angle
				if nowX > c.x {
					c.x = 0
				} else {
					c.x -= nowX
				}
			}
			if Angle < 0 {
				nowX := -1 * Angle
				if nowX > c.room.Size_X-c.x {
					c.x = c.room.Size_X
				} else {
					c.x += nowX
				}
			}
		}
		select {
		case <-t2.C:
			if c.x < 3000 && c.y < 2100 {
				forSend := make([]byte, 9)
				forSend[0] = 7
				a := make([]byte, 4)
				binary.LittleEndian.PutUint32(a, c.room.Base_A.HP)
				forSend[1] = a[3]
				forSend[2] = a[2]
				forSend[3] = a[1]
				forSend[4] = a[0]
				a = make([]byte, 4)
				binary.LittleEndian.PutUint32(a, c.room.Base_A.Shield)
				forSend[5] = a[3]
				forSend[6] = a[2]
				forSend[7] = a[1]
				forSend[8] = a[0]
				c.send <- forSend
			}
			if c.x > 13000 && c.y > 13900 {
				forSend := make([]byte, 9)
				forSend[0] = 8
				a := make([]byte, 4)
				binary.LittleEndian.PutUint32(a, c.room.Base_B.HP)
				forSend[1] = a[3]
				forSend[2] = a[2]
				forSend[3] = a[1]
				forSend[4] = a[0]
				a = make([]byte, 4)
				binary.LittleEndian.PutUint32(a, c.room.Base_B.Shield)
				forSend[5] = a[3]
				forSend[6] = a[2]
				forSend[7] = a[1]
				forSend[8] = a[0]
				c.send <- forSend
			}
			for q, u := range c.room.Asteroids {
				if u == nil || u.hide {
					continue
				}
				if u.x >= c.x {
					if u.y >= c.y {
						if u.y-c.y > 1350 || u.x-c.x > 2600 {
							continue
						}
						forSend := make([]byte, 19)
						forSend[0] = 9
						forSend[1] = q
						forSend[2] = u.style
						f := Float64bytes(u.x)
						forSend[3] = f[0]
						forSend[4] = f[1]
						forSend[5] = f[2]
						forSend[6] = f[3]
						forSend[7] = f[4]
						forSend[8] = f[5]
						forSend[9] = f[6]
						forSend[10] = f[7]
						f = Float64bytes(u.y)
						forSend[11] = f[0]
						forSend[12] = f[1]
						forSend[13] = f[2]
						forSend[14] = f[3]
						forSend[15] = f[4]
						forSend[16] = f[5]
						forSend[17] = f[6]
						forSend[18] = f[7]
						c.send <- forSend
						continue
					}
					if c.y > u.y {
						if c.y-u.y > 1350 || u.x-c.x > 2600 {
							continue
						}
						forSend := make([]byte, 19)
						forSend[0] = 9
						forSend[1] = q
						forSend[2] = u.style
						f := Float64bytes(u.x)
						forSend[3] = f[0]
						forSend[4] = f[1]
						forSend[5] = f[2]
						forSend[6] = f[3]
						forSend[7] = f[4]
						forSend[8] = f[5]
						forSend[9] = f[6]
						forSend[10] = f[7]
						f = Float64bytes(u.y)
						forSend[11] = f[0]
						forSend[12] = f[1]
						forSend[13] = f[2]
						forSend[14] = f[3]
						forSend[15] = f[4]
						forSend[16] = f[5]
						forSend[17] = f[6]
						forSend[18] = f[7]
						c.send <- forSend
						continue
					}
				}
				if c.x > u.x {
					if u.y >= c.y {
						if u.y-c.y > 1350 || c.x-u.x > 2600 {
							continue
						}
						forSend := make([]byte, 19)
						forSend[0] = 9
						forSend[1] = q
						forSend[2] = u.style
						f := Float64bytes(u.x)
						forSend[3] = f[0]
						forSend[4] = f[1]
						forSend[5] = f[2]
						forSend[6] = f[3]
						forSend[7] = f[4]
						forSend[8] = f[5]
						forSend[9] = f[6]
						forSend[10] = f[7]
						f = Float64bytes(u.y)
						forSend[11] = f[0]
						forSend[12] = f[1]
						forSend[13] = f[2]
						forSend[14] = f[3]
						forSend[15] = f[4]
						forSend[16] = f[5]
						forSend[17] = f[6]
						forSend[18] = f[7]
						c.send <- forSend
						continue
					}
					if c.y > u.y {
						if c.y-u.y > 1350 || c.x-u.x > 2600 {
							continue
						}
						forSend := make([]byte, 19)
						forSend[0] = 9
						forSend[1] = q
						forSend[2] = u.style
						f := Float64bytes(u.x)
						forSend[3] = f[0]
						forSend[4] = f[1]
						forSend[5] = f[2]
						forSend[6] = f[3]
						forSend[7] = f[4]
						forSend[8] = f[5]
						forSend[9] = f[6]
						forSend[10] = f[7]
						f = Float64bytes(u.y)
						forSend[11] = f[0]
						forSend[12] = f[1]
						forSend[13] = f[2]
						forSend[14] = f[3]
						forSend[15] = f[4]
						forSend[16] = f[5]
						forSend[17] = f[6]
						forSend[18] = f[7]
						c.send <- forSend
						continue
					}
				}
			}
			CurrRoom.RW.RLock()
			for _, u := range CurrRoom.Clients {
				if u == nil {
					continue
				}
				if u.x >= c.x {
					if u.y >= c.y {
						if u.y-c.y > 1350 || u.x-c.x > 2600 {
							continue
						}
						c.send <- GetByte(u)
						continue
					}
					if c.y > u.y {
						if c.y-u.y > 1350 || u.x-c.x > 2600 {
							continue
						}
						c.send <- GetByte(u)
						continue
					}
				}
				if c.x > u.x {
					if u.y >= c.y {
						if u.y-c.y > 1350 || c.x-u.x > 2600 {
							continue
						}
						c.send <- GetByte(u)
						continue
					}
					if c.y > u.y {
						if c.y-u.y > 1350 || c.x-u.x > 2600 {
							continue
						}
						c.send <- GetByte(u)
						continue
					}
				}
			}
			CurrRoom.RW.RUnlock()
			c.room.ShellsMut.RLock()
			for q, u := range c.room.Shells {
				if u.x >= c.x {
					if u.y >= c.y {
						if u.y-c.y > 1300 || u.x-c.x > 2000 {
							continue
						}
						c.send <- GetByteShell(u, q)
						continue
					}
					if c.y > u.y {
						if c.y-u.y > 1300 || u.x-c.x > 2000 {
							continue
						}
						c.send <- GetByteShell(u, q)
						continue
					}
				}
				if c.x > u.x {
					if u.y >= c.y {
						if u.y-c.y > 1300 || c.x-u.x > 2000 {
							continue
						}
						c.send <- GetByteShell(u, q)
						continue
					}
					if c.y > u.y {
						if c.y-u.y > 1300 || c.x-u.x > 2000 {
							continue
						}
						c.send <- GetByteShell(u, q)
						continue
					}
				}
			}
			c.room.ShellsMut.RUnlock()
			c.room.EatMut.RLock()
			for q, u := range c.room.Eats {
				if u.x >= c.x {
					if u.y >= c.y {
						if u.y-c.y > 1200 || u.x-c.x > 1950 {
							continue
						}
						c.send <- GetByteEat(u, q)
						continue
					}
					if c.y > u.y {
						if c.y-u.y > 1200 || u.x-c.x > 1950 {
							continue
						}
						c.send <- GetByteEat(u, q)
						continue
					}
				}
				if c.x > u.x {
					if u.y >= c.y {
						if u.y-c.y > 1200 || c.x-u.x > 1950 {
							continue
						}
						c.send <- GetByteEat(u, q)
						continue
					}
					if c.y > u.y {
						if c.y-u.y > 1200 || c.x-u.x > 1950 {
							continue
						}
						c.send <- GetByteEat(u, q)
						continue
					}
				}
			}
			c.room.EatMut.RUnlock()
			break
		default:
			break
		}
	}
}
