package websocket

import (
	"fmt"

	"../gamemodel"
)

/*
 * Game
 *  - Register: 	Channel for clients to register to the game
 *  - Unregister:	Channel for clients to unregister from the game
 *  - Clients:		Map of clients in the game
 *  - Broadcast:	Channel for clients to broadcast messages to neighbors in game
 *  - Model:			Model of the current game
 */
type Game struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
	Model      *gamemodel.GameModel
}

// NewGame function returns a game model
func NewGame() *Game {
	return &Game{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
		Model:      gamemodel.NewGameModel(),
	}
}

// Start function listens to different channels
func (game *Game) Start() {
	for {
		select {
		case client := <-game.Register:
			game.Clients[client] = true
			fmt.Println("Size of Connection game: ", len(game.Clients))
			for client := range game.Clients {
				fmt.Println(client)
				client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
			}
			break
		case client := <-game.Unregister:
			delete(game.Clients, client)
			fmt.Println("Size of Connection Game: ", len(game.Clients))
			for client := range game.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
			}
			break
		case message := <-game.Broadcast:
			fmt.Println("Sending message to all clients in Game")
			for client := range game.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
