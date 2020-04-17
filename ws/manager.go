package ws

import (
	"encoding/json"
	"log"
)

var Handlers = make(map[string]WsHandler)


type WsHandler func(c *Context)

type Context struct {
	CMD string `json:"cmd"`
	Data interface{} `json:"data"`
	Client *Client `json:"-"`
}


func AddHandler(cmd string,handler WsHandler) {
	Handlers[cmd] = handler
}

func (c *Context) Send(data interface{}) {
	bytes,_ := json.Marshal(data)
	c.Client.send <- bytes
}
func (c *Context) SendString(data string) {
	log.Println(data)
	c.Client.send <- []byte(data)
}

//对房间广播消息
func (c *Context) SendRoomBroadcast(data RoomBroadcast) {
	c.Client.RoomBroadcast(data)
}


func invoke(c *Context) {
	if f,ok:=Handlers[c.CMD];ok {
		f(c)
	}
}