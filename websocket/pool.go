package websocket

import "fmt"

/**
* Pool
*  - Register: Channel for clients to register to the pool
*  - Unregister: Channel for clients to unregister from the pool
*  - Clients: A map of clients connected to the pool
*  - Broadcast: Channel for messaging all clients in pool
*  - Neighbors: Potentially useless :)
 */
type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

/*
* NewPool
* @return a generated pool
 */
func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))

			//for client := range pool.Clients {
			//	fmt.Println(client)
			//	client.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined..."})
			//}
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected..."})
			}
			break
		case message := <-pool.Broadcast:
			//fmt.Println("The message is: ", message)
			fmt.Println("Sending message to all clients in Pool")
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}

	}
}
