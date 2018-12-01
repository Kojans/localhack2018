package main

import "websocket-master"

type Client struct {
    ws   *websocket.Conn
    send chan []byte

    player *Player
}

type Player struct {
    Ship       int
    MainGuns   map[int]int
    AddGuns    map[int]int
    Generators map[int]int
    Engines    map[int]int

    HitPoints int
    Shield    int
    ShieldRP  int

    Speed float64

    x         float64
    y         float64
    base      int
    points    int
    resources [3]int
}
