package main

func CreateEat(x float64, y float64, room *Room, points uint32) bool {
	room.EatMut.RLock()
	for _, i := range room.Eats {
		if i.x > x {
			if i.y > y {
				if (i.x-x)+(i.y-y) < 60 {
					room.EatMut.RUnlock()
					return false
				}
			} else {
				if (i.x-x)+(y-i.y) < 60 {
					room.EatMut.RUnlock()
					return false
				}
			}
		} else {
			if i.y > y {
				if (x-i.x)+(i.y-y) < 60 {
					room.EatMut.RUnlock()
					return false
				}
			} else {
				if (x-i.x)+(y-i.y) < 60 {
					room.EatMut.RUnlock()
					return false
				}
			}
		}
	}
	room.EatMut.RUnlock()
	room.EatNMut.Lock()
	room.EatMut.Lock()
	room.Eats[room.EatN] = &Eat{x, y, points}
	room.EatN += 1
	room.EatNMut.Unlock()
	room.EatMut.Unlock()
	return true
}

type Eat struct {
	x, y   float64
	t uint32
}
