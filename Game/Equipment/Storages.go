package Equipment

var Storages = []*Storage{
    &Storage {
        50,
        500,
    },
    &Storage {
        100,
        1000,
    },
    &Storage {
        200,
        2000,
    },
}

type Storage struct {
    Space int
    Cost int
}