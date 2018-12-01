package main

type Generator struct {
	t uint8
}

type Accelerator struct {
	t uint8
}

var Engins = []*Engine{
	NewEngine(0),
	NewEngine(1),
	NewEngine(2),
}

func NewEngine(t uint8) *Engine {
	if t == 0 {
		return &Engine{
			0,
			100,
			3,
			5,
			2,
		}
	}
	if t == 1 {
		return &Engine{
			1,
			250,
			5,
			9,
			4,
		}
	}
	if t == 2 {
		return &Engine{
			2,
			500,
			7,
			12,
			4,
		}
	}
	return nil
}

type Engine struct {
	t uint8
	Energy   uint16
	Generate uint8
	Speed    float64
	Use      uint8
}
