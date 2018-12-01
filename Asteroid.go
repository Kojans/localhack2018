package main

import (
	"time"
	"math/rand"
	"math"
)

func (a *Asteroid) start() {
	add := float64(275)
	if a.style == 1 {
		add = 375
	}
	if a.style == 2 {
		add = 475
	}
	for {
		b := rand.Int31n(a.max_spawn_time) + a.min_spawn_time
		<-time.NewTicker(time.Duration(b) * time.Second).C
		b = rand.Int31n(a.min_point) + a.max_point
		for ; b > 0; b-- {
			for {
				rand.Seed(time.Now().Unix()*1000 + time.Now().UnixNano()/500)
				mul := float64(rand.Int31n(200))
				rand.Seed(time.Now().UnixNano()/500 + time.Now().Unix()*800)
				Angle := rand.Float64() * 2 * math.Pi
				x := (mul + add) * math.Cos(Angle)+a.x
				y := (mul + add) * math.Sin(Angle)+a.y

				if CreateEat(x, y, a.room, a.points) {
					break
				}
			}
		}
	}
}

func NewAsteroid(t float64, x, y float64, room *Room) *Asteroid {
	switch t {
	case 0:
		return &Asteroid{
			x,
			y,
			false,
			0,
			100,
			1000,
			1,
			5,
			15,
			30,
			40,
			60,
			room,
			0,
		}
	case 1:
		return &Asteroid{
			x,
			y,
			false,
			1,
			150,
			3000,
			1,
			3,
			30,
			60,
			60,
			120,
			room,
			1,
		}
	case 2:
		return &Asteroid{
			x,
			y,
			false,
			2,
			100,
			1000,
			0,
			0,
			0,
			0,
			120,
			180,
			room,
			2,
		}
	}
	return nil
}

type Asteroid struct {
	x                float64
	y                float64
	hide             bool
	style            uint8
	weight           uint16
	size             uint16
	min_point        int32
	max_point        int32
	min_spawn_time   int32
	max_spawn_time   int32
	min_respawn_time int32
	max_respawn_time int32
	room             *Room
	points           uint32
}
