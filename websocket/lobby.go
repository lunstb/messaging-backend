package websocket

import (
	"encoding/json"
	"fmt"
)

type NeighborMessage struct {
	Clients map[*Client]bool
	Message string
}

/*
 * Lobby
 *  - Clients: 		Map of clients in the lobby
 *  - Message: 		What is this for??
 *  - Register: 	Channel for clients to register for the lobby
 *  - Unregister:	Channel for clients to unregister from the lobby
 *  - Broadcast:	Channel for clients to broadcast messages to clients
 */
type Lobby struct {
	Clients    map[*Client]bool
	Message    string
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
}

// NewLobby returns a new lobby
func NewLobby() *Lobby {
	return &Lobby{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

// Start function sets up listener for channels
func (lobby *Lobby) Start() {
	for {
		select {
		case client := <-lobby.Register:
			lobby.Clients[client] = true
			fmt.Println("Size of Connection Lobby: ", len(lobby.Clients))

			responseJSON := Response{Pool: "Lobby"}
			messageJSON := MessageToClient{Type: "ChangedPool", Response: &responseJSON}
			messageByte, _ := json.Marshal(messageJSON)

			fmt.Println(string(messageByte))

			client.Conn.WriteJSON(Message{Type: 1, Body: string(messageByte)})
			for client := range lobby.Clients {
				fmt.Println(client)
				client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
			}
			break
		case client := <-lobby.Unregister:
			delete(lobby.Clients, client)
			fmt.Println("Size of Connection Lobby: ", len(lobby.Clients))
			for client := range lobby.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
			}
			break
		case message := <-lobby.Broadcast:
			fmt.Println("Sending message to all clients in Lobby")
			for client := range lobby.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
