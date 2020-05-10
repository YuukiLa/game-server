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

	roomBroadcast    chan RoomBroadcast
	roomBroadcastAll chan RoomBroadcast
	roomSendUser     chan RoomBroadcast
	rooms            map[string][]*Client
}

func NewHub() *Hub {
	return &Hub{
		register:         make(chan *Client),
		unregister:       make(chan *Client),
		enterRoom:        make(chan *Client),
		exitRoom:         make(chan *Client),
		roomBroadcast:    make(chan RoomBroadcast),
		roomBroadcastAll: make(chan RoomBroadcast),
		roomSendUser:     make(chan RoomBroadcast),
		clients:          make(map[string]*Client),
		rooms:            make(map[string][]*Client),
	}
}

func (h *Hub) Run() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("hub error", err)
		}
	}()
	for {
		select {
		case client := <-h.register:
			log.Println("链接")
			if c, ok := h.clients[client.ID]; ok {
				log.Println("有未断开的连接")
				h.unregister <- c
			}
			h.clients[client.ID] = client

		case client := <-h.unregister:
			log.Printf("client: %+v", client)
			log.Printf("clients:%+v,%+v", h.clients[client.ID])
			if c, ok := h.clients[client.ID]; ok {
				log.Println("关闭1")
				if c.RoomID != "" {
					log.Println("关闭2")
					go func() {
						h.exitRoom <- c
					}()
				}
				log.Println("关闭3")
				delete(h.clients, c.ID)
				log.Println("关闭4")
				close(c.send)
				log.Println("关闭5")
				log.Printf("关闭后的clients:%+v", h.clients)
			}
		case client := <-h.enterRoom:
			log.Println("进入房间")
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
		case roomBroadcast := <-h.roomSendUser:
			log.Println("给指定用户发送")
			if clients, ok := h.rooms[roomBroadcast.RoomID]; ok {
				for _, client := range clients {
					if client.ID == roomBroadcast.ClientID {
						data, _ := json.Marshal(roomBroadcast.Data)
						client.send <- data
					}
				}
			}
		case roomBroadcast := <-h.roomBroadcastAll:
			log.Println("广播所有")
			log.Println(roomBroadcast.RoomID)
			log.Printf("%+v", h.rooms)
			if clients, ok := h.rooms[roomBroadcast.RoomID]; ok {
				log.Printf("%+v", ok)
				for _, client := range clients {
					data, _ := json.Marshal(roomBroadcast.Data)
					client.send <- data

				}
			}
		}
	}
}
