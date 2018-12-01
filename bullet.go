package main

var Bullets = make(map[*Bullet]struct{})

type Bullet struct {
	x float64
	y float64
}
