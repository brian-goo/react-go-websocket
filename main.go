package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Stream struct {
	Date  string  `json:"date"`
	Time  string  `json:"time"`
	Price float64 `json:"price"`
}

func getNext() func() float64 {
	v := 1.0
	return func() float64 {
		rand.Seed(time.Now().UnixNano())
		r := math.Max(rand.NormFloat64(), -0.1)
		t := math.Min(r, 0.1)
		v = v * (1 + t)
		if v < 0.2 {
			v = 0.2
		}
		return math.Round(v*100) / 100
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	defer c.Close()

	next := getNext()

	for i := 0; i < 10; i++ {
		st := Stream{
			time.Now().Format("2006-01-02"),
			time.Now().Format("15:04:05"),
			// float64(i)*0.01 + 2,
			next(),
		}
		err := c.WriteJSON(st)
		if err != nil {
			fmt.Println("err:", err)
			break
		} else {
			fmt.Println(st)
		}
		time.Sleep(time.Second)
	}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./js")))
	http.HandleFunc("/ws", echo)

	fmt.Println("server starting...")
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
