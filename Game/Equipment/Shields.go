package Equipment

import "time"

var Shields = []*Shield{
    &Shield{
        2,
        5,
        10,
        3,
        5000,
    },
    &Shield{
        3,
        10,
        50,
        3,
        12000,
    },
}

type Shield struct {
    Append       float64
    Recovery     int
    SpendEnergyP int
    Time         time.Duration
    Cost         int
}
