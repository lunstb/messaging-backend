package main

import (
	"fmt"
	"log"
	"net/http"

	"example.com/main/stringgen"
	"example.com/main/websocket"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func testPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Test is successful")
	fmt.Println("Endpoint Hit: test page")
}

func serveWs(pool *websocket.Pool, lobby *websocket.Lobby, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: WebSocket")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		ID:    stringgen.String(10),
		Conn:  conn,
		Pool:  pool,
		Lobby: lobby,
	}

	pool.Register <- client
	lobby.Register <- client
	client.Read()
}

func setupRoutes() {
	websocket.Test()

	pool := websocket.NewPool()
	lobby := websocket.NewLobby()

	go pool.Start()
	go lobby.Start()

	http.HandleFunc("/", homePage)
	http.HandleFunc("/test", testPage)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, lobby, w, r)
	})
}

func main() {
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
