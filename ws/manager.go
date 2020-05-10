package ws

import (
	"encoding/json"
	"log"
)

var Handlers = make(map[string]WsHandler)

type WsHandler func(c *Context)

type Context struct {
	CMD    string      `json:"cmd"`
	Data   interface{} `json:"data"`
	Client *Client     `json:"-"`
}

func AddHandler(cmd string, handler WsHandler) {
	Handlers[cmd] = handler
}

func (c *Context) Bind(model interface{}) error {
	bytes, _ := json.Marshal(c.Data)
	return json.Unmarshal(bytes, model)
}

func (c *Context) Send(data interface{}) {
	bytes, _ := json.Marshal(data)
	c.Client.send <- bytes
}
func (c *Context) SendString(data string) {
	log.Println(data)
	c.Client.send <- []byte(data)
}

//对房间广播消息(除了自己)
func (c *Context) SendRoomBroadcast(data RoomBroadcast) {
	c.Client.RoomBroadcast(data)
}

//对房间广播消息(包括自己)
func (c *Context) SendRoomBroadcastAll(data RoomBroadcast) {
	c.Client.RoomBroadcastAll(data)
}

//对房间具体用户下发信息
func (c *Context) SendRoomUser(data RoomBroadcast) {
	c.Client.RoomSendUser(data)
}

func invoke(c *Context) {
	if f, ok := Handlers[c.CMD]; ok {
		f(c)
	}
}
