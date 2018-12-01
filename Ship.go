package main

var Ships = []*Ship{
	NewShip(0),
	NewShip(1),
	NewShip(2),
}

func NewShip(t uint8) *Ship {
	switch t {
	case 0:
		return &Ship{
			0,
			1,
			2,
			1,
			2,
			2,
			150,
			10,
			1,
			1.5,
			50,
		}
	case 1:
		return &Ship{
			1,
			1,
			3,
			2,
			2,
			1,
			350,
			10,
			1,
			1,
			100,
		}
	case 2:
		return &Ship{
			2,
			2,
			6,
			3,
			4,
			2,
			1500,
			10,
			1,
			0.3,
			300,
		}
	}
	return nil
}

type Ship struct {
	Type            uint8
	Main_Guns       uint8
	Additional_Guns uint8
	Generators      uint8
	Accelerators    uint8
	Engine          uint8
	HP              uint16
	Shield          uint8
	Shield_recovery uint8
	Speed           float64
	Space           uint16
}
