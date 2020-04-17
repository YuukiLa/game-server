package ws

import (
	"encoding/json"
	"log"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[string]*Client
	// Register requests from the clients.
	register chan *Client

	enterRoom chan *Client

	exitRoom chan *Client
	// Unregister requests from clients.
	unregister chan *Client

	roomBroadcast chan RoomBroadcast

	rooms map[string][]*Client
}

func NewHub() *Hub {
	return &Hub{
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		enterRoom:     make(chan *Client),
		exitRoom:      make(chan *Client),
		roomBroadcast: make(chan RoomBroadcast),
		clients:       make(map[string]*Client),
		rooms:         make(map[string][]*Client),
	}
}

func (h *Hub) Run() {
	defer func() {
		if err := recover();err!=nil{
			log.Println("hub error",err)
		}
	}()
	for {
		select {
		case client := <-h.register:
			h.clients[client.ID] = client
		case client := <-h.unregister:
			if _, ok := h.clients[client.ID]; ok {
				if client.RoomID!="" {
					h.exitRoom <- client
				}
				delete(h.clients, client.ID)
				close(client.send)
			}
		case client := <-h.enterRoom:
			h.rooms[client.RoomID] = append(h.rooms[client.RoomID], client)
		case client := <-h.exitRoom:
			for i, c := range h.rooms[client.RoomID] {
				if c.ID == client.ID {
					h.rooms[client.RoomID] = append(h.rooms[client.RoomID][:i], h.rooms[client.RoomID][i+1:]...)
					client.RoomID = ""
					continue
				}
			}
		case roomBroadcast := <-h.roomBroadcast:
			if clients, ok := h.rooms[roomBroadcast.RoomID]; ok {
				for _, client := range clients {
					if client.ID != roomBroadcast.ClientID {
						data, _ := json.Marshal(roomBroadcast.Data)
						client.send <- data
					}
				}
			}
		}
	}
}
