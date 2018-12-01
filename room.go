package main

import "sync"

func NewRoom(x, y float64) *Room {
	return &Room{
		make(map[*Client]bool),
		sync.RWMutex{},
		make(map[*Client]bool),
		sync.RWMutex{},
		x,
		y,
		NewBase(),
		NewBase(),
		0,
		sync.Mutex{},
		make(map[uint32]*Shell),
		sync.RWMutex{},
		make(map[uint32]*Eat),
		sync.RWMutex{},
		0,
		sync.Mutex{},
		make(map[uint8]*Asteroid),
		sync.RWMutex{},
		make(map[string]*Client),
		sync.RWMutex{},
	}
}

type Room struct {
	Command_A map[*Client]bool
	AMut      sync.RWMutex

	Command_B map[*Client]bool
	BMut      sync.RWMutex

	Size_X float64
	Size_Y float64

	Base_A *Base
	Base_B *Base

	ShellId  uint32
	ShellMut sync.Mutex

	Shells    map[uint32]*Shell
	ShellsMut sync.RWMutex

	Eats   map[uint32]*Eat
	EatMut sync.RWMutex

	EatN    uint32
	EatNMut sync.Mutex

	Asteroids    map[uint8]*Asteroid
	AsteroidsMut sync.RWMutex

	Clients map[string]*Client
	RW      sync.RWMutex
}
