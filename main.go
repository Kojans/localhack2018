package main

import (
	"log"
	"net/http"
	"websocket-master"
	"strconv"
)

var CurrRoom *Room

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HTML_MAIN(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./home.html")
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		return
	}
	CurrRoom.RW.RLock()
	conn, err := upgrader.Upgrade(w, r, nil)
	if len(CurrRoom.Clients) > 250 {
		CurrRoom.RW.RUnlock()
		//conn.WriteMessage(websocket.BinaryMessage, []byte("1"))
		w.Write([]byte("Server is full"))
		return
	}
	CurrRoom.RW.RUnlock()
	if err != nil {
		return
	}
	CurrRoom.RW.Lock()
	if _, ok := CurrRoom.Clients[name]; !ok {
		CurrRoom.Clients[name] = nil
	} else {
		CurrRoom.RW.Unlock()
		return
	}
	CurrRoom.RW.Unlock()
	CurrRoom.AMut.RLock()
	CurrRoom.BMut.RLock()
	x := float64(16000)
	y := float64(16000)
	command := false
	if len(CurrRoom.Command_A) < len(CurrRoom.Command_B){
		x = 0
		y = 0
		command = true
	}
	CurrRoom.AMut.RUnlock()
	CurrRoom.BMut.RUnlock()

	c := NewClient(conn, name, x, y, CurrRoom, command)
	c.Engine[0] = Engins[0]
	c.SetSpeed()
	if command {
		CurrRoom.AMut.Lock()
		CurrRoom.Command_A[c] = false
		CurrRoom.AMut.Unlock()
	} else {
		CurrRoom.BMut.Lock()
		CurrRoom.Command_B[c] = false
		CurrRoom.BMut.Unlock()
	}

	CurrRoom.RW.Lock()
	CurrRoom.Clients[name] = c
	CurrRoom.RW.Unlock()

	go c.Shooting()
	go c.writer()
	go c.reader()
	go c.RecoveryShield()
	go c.start()
	c.collision()
}

func HTML_CODE(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Code Page</title>
</head>
<body>
<form action="#">
    <input name="l">
    <input name="x">
    <input name="y">
    <input type="submit">
</form>
</body>
	`))
	if r.FormValue("l") != "" && r.FormValue("x") != "" && r.FormValue("y") != "" {
		CurrRoom.RW.RLock()
		if q, ok := CurrRoom.Clients[r.FormValue("l")]; ok {
			x, _ := strconv.Atoi(r.FormValue("x"))
			y, _ := strconv.Atoi(r.FormValue("y"))
			q.x = float64(x)
			q.y = float64(y)
		}
		CurrRoom.RW.RUnlock()
	}
}

func main() {
	CurrRoom = NewRoom(16000, 16000)

	CurrRoom.Asteroids[0] = NewAsteroid(0, 550, 1550, CurrRoom)
	CurrRoom.Asteroids[1] = NewAsteroid(0, 1550, 550, CurrRoom)
	CurrRoom.Asteroids[2] = NewAsteroid(0, 2350, 2350, CurrRoom)

	CurrRoom.Asteroids[3] = NewAsteroid(0, 16000-550, 16000-1550, CurrRoom)
	CurrRoom.Asteroids[4] = NewAsteroid(0, 16000-1550, 16000-550, CurrRoom)
	CurrRoom.Asteroids[5] = NewAsteroid(0, 16000-2350, 16000-2350, CurrRoom)

	go CurrRoom.Asteroids[0].start()
	go CurrRoom.Asteroids[1].start()
	go CurrRoom.Asteroids[2].start()
	go CurrRoom.Asteroids[3].start()
	go CurrRoom.Asteroids[4].start()
	go CurrRoom.Asteroids[5].start()


	CurrRoom.Asteroids[6] = NewAsteroid(1, 4650, 2150, CurrRoom)
	CurrRoom.Asteroids[7] = NewAsteroid(1, 2150, 4650, CurrRoom)
	CurrRoom.Asteroids[8] = NewAsteroid(1, 4650, 4650, CurrRoom)

	CurrRoom.Asteroids[9] = NewAsteroid(1, 16000-4650, 16000-2150, CurrRoom)
	CurrRoom.Asteroids[10] = NewAsteroid(1, 16000-2150, 16000-4650, CurrRoom)
	CurrRoom.Asteroids[11] = NewAsteroid(1, 16000-4650, 16000-4650, CurrRoom)

	go CurrRoom.Asteroids[6].start()
	go CurrRoom.Asteroids[7].start()
	go CurrRoom.Asteroids[8].start()
	go CurrRoom.Asteroids[9].start()
	go CurrRoom.Asteroids[10].start()
	go CurrRoom.Asteroids[11].start()


	CurrRoom.Asteroids[12] = NewAsteroid(2, 8000, 8000, CurrRoom)
	CurrRoom.Asteroids[13] = NewAsteroid(2, 6500, 6500, CurrRoom)
	CurrRoom.Asteroids[14] = NewAsteroid(2, 6500, 16000-6500, CurrRoom)
	CurrRoom.Asteroids[15] = NewAsteroid(2, 16000-6500, 16000-6500, CurrRoom)
	CurrRoom.Asteroids[16] = NewAsteroid(2, 16000-6500, 6500, CurrRoom)

	log.Println("Start")

	log.Println("23.08.2017 build")

	http.HandleFunc("/", HTML_MAIN)
	http.HandleFunc("/game", HTML_MAIN)
	http.HandleFunc("/code", HTML_CODE)
	http.HandleFunc("/ws", serveWs)
	http.ListenAndServe(":5692", nil)
}
