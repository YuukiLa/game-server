package handler

import (
	"encoding/json"
	"github/com/yuuki80code/game-server/cmd_router/handler/model"
	"github/com/yuuki80code/game-server/controller/response"
	"github/com/yuuki80code/game-server/redis"
	"github/com/yuuki80code/game-server/ws"
	"log"
)

func EnterRoom(c *ws.Context) {
	cache := redis.Cache
	var roomId model.RoomID
	err := c.Bind(&roomId)
	log.Println("invoke")
	if err != nil {
		c.SendString(err.Error())
		return
	}
	if err != nil {
		c.SendString(err.Error())
		return
	}
	users := make([]response.UserInfoResp, 0)
	data, err := cache.HGetAll(roomId.RoomID).Result()
	if err != nil {
		log.Println(err)
		c.SendString(err.Error())
		return
	}

	for _, value := range data {
		var user response.UserInfoResp
		json.Unmarshal([]byte(value), &user)
		users = append(users, user)
	}
	if c.Client.RoomID == "" {
		c.Client.RoomID = roomId.RoomID
		c.Client.EnterRoom()
	}
	log.Println(users)
	result := ws.Result{
		CMD:  c.CMD,
		Data: users,
		Msg:  "",
	}
	broadcast := ws.RoomBroadcast{
		RoomID:   roomId.RoomID,
		Data:     result,
		ClientID: "",
	}
	c.SendRoomBroadcastAll(broadcast)
}
