package main

import "time"

var Guns = []*Gun{
	&Gun{2*time.Second,10,10},
	&Gun{time.Second*3/2,30,15},
	&Gun{time.Second,50,25},
}

type Gun struct {
	time.Duration
	Damage      uint8
	Consumption uint8
}

type GunOnShip struct{
	x,y,t uint8
}