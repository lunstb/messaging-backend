package main

import (
	"fmt"
	"log"
	"net/http"

	"example.com/main/database"
	"example.com/main/stringgen"
	"example.com/main/websocket"
)

var db *database.Database

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func testPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test is successful")
	fmt.Println("Endpoint Hit: test page")
}

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: WebSocket")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		ID:   stringgen.String(10),
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func setupRoutes() {
	db = database.Connect()

	db.Test()

	pool := websocket.NewPool()

	go pool.Start()

	http.HandleFunc("/", homePage)
	http.HandleFunc("/test", testPage)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})
}

func main() {
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
