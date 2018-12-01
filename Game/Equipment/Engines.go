package Equipment

import "time"

var Engines = []*Engine{
    &Engine{
        100,
        3,
        time.Second,
        5,
        2,
        2,
        0,
    },
    &Engine{
        250,
        5,
        time.Second,
        9,
        3,
        4,
        5000,
    },
    &Engine{
        500,
        7,
        time.Second,
        12,
        4,
        5,
        8000,
    },
}

type Engine struct {
    Energy      int
    GenerateP   int
    GenerateT   time.Duration
    SpeedMax    int
    SpeedMin    int
    SpendEnergy int
    Cost        int
}
