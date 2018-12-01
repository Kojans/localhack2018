package Equipment

var Generators = []*Generator{
    &Generator{
        500,
        2,
        3000,
    },
    &Generator{
        1000,
        3,
        8000,
    },
}

type Generator struct {
    Energy   int
    Generate int
    Cost     int
}
