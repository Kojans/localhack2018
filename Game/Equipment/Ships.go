package Equipment

import "time"

var Ships = []*Ship{
    &Ship{
        [][2]float64{
            {0, 0},
        },
        [][2]float64{
            {1, 0},
            {2, 0},
        },
        1,
        2,
        1,
        50,
        150,
        10,
        1,
        3*time.Second,
        1.5,
        0,
    },
    &Ship{
        [][2]float64{
            {0, 0},
        },
        [][2]float64{
            {1, 0},
            {2, 0},
            {3, 0},
        },
        2,
        2,
        1,
        100,
        350,
        50,
        1,
        3*time.Second,
        1,
        5000,
    },
    &Ship{
        [][2]float64{
            {0, 0},
            {1, 0},
        },
        [][2]float64{
            {1, 0},
            {2, 0},
            {3, 0},
            {4, 0},
            {5, 0},
            {6, 0},
        },
        3,
        4,
        2,
        300,
        1500,
        200,
        1,
        3*time.Second,
        0.3,
        25000,
    },
}

type Ship struct {
    MainGuns [][2]float64
    AddGuns  [][2]float64

    Generators   int
    Accelerators int
    Engines      int
    Space        int

    HitPoints    int
    ShieldPoints int

    ShieldRP int
    ShieldRT time.Duration

    Speed float64
    Cost  int
}
