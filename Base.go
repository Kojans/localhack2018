package main

import "time"

func (b *Base) Start() {
	for {
		<-b.T.C
		if b.Max_Shield == b.Shield {
			continue
		}
		if b.Max_Shield > b.Shield+b.Per_Tick {
			b.Shield += b.Per_Tick
			continue
		}
		b.Shield = b.Max_Shield
	}
}

func NewBase() *Base {
	return &Base{
		0,
		10000,
		100,
		1,
		100,
		true,
		10,
		time.NewTicker(3 * time.Second),
	}
}

type Base struct {
	Points     uint32
	HP         uint32
	Shield     uint32
	Max_Shield uint32
	Per_Tick   uint32
	IsHide     bool
	Damage     uint16
	T          *time.Ticker
}
