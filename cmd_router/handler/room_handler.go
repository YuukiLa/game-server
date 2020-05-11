package handler

import (
	"encoding/json"
	"github/com/yuuki80code/game-server/cmd_router/handler/model"
	"github/com/yuuki80code/game-server/config/define"
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
	data, err := cache.LRange(roomId.RoomID, 0, -1).Result()
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
	var roomConfig define.RoomConfig
	roomConfigStr, _ := cache.HGet(define.RedisRoomConfigKey, roomId.RoomID).Result()
	json.Unmarshal([]byte(roomConfigStr), &roomConfig)
	if c.Client.ID == roomConfig.Master {
		c.Send(ws.Result{
			CMD:  "10002",
			Data: make(map[string]interface{}),
			Msg:  "你是房主",
		})
	}
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

func StartGame(c *ws.Context) {
	cache := redis.Cache
	var roomConfig define.RoomConfig
	roomId := c.Client.RoomID
	if roomId == "" {

	}
	roomConfigStr, _ := cache.HGet(define.RedisRoomConfigKey, roomId).Result()
	json.Unmarshal([]byte(roomConfigStr), &roomConfig)
	if c.Client.ID == roomConfig.Master {
		roomConfig.Status = 1
		roomConfig.CurrUser = c.Client.ID
		roomConfig.CurrWord = "马"
		roomConfig.Round = 1
		cache.HSet(define.RedisRoomConfigKey, roomId, roomConfig)
		result := ws.Result{
			CMD:  "10010",
			Data: roomConfig,
			Msg:  "",
		}
		broadcast := ws.RoomBroadcast{
			RoomID:   roomId,
			Data:     result,
			ClientID: "",
		}
		c.SendRoomBroadcastAll(broadcast)
	}
}

func GameParam(c *ws.Context) {
	cache := redis.Cache
	var roomConfig define.RoomConfig

	roomId := c.Client.RoomID
	if roomId == "" {

	}
	roomConfigStr, _ := cache.HGet(define.RedisRoomConfigKey, roomId).Result()
	json.Unmarshal([]byte(roomConfigStr), &roomConfig)
	//len, err := cache.LLen(roomId).Result()
	//if err != nil {
	//
	//}
	//userStr, _ := cache.LIndex(roomId, int64(roomConfig.Round-1)%len).Result()
	//var user response.UserInfoResp
	//json.Unmarshal([]byte(userStr), &user)
	if c.Client.ID == roomConfig.CurrUser {
		result := ws.Result{
			CMD:  "11000",
			Data: roomConfig,
			Msg:  "",
		}
		broadcast := ws.RoomBroadcast{
			RoomID:   roomId,
			Data:     result,
			ClientID: c.Client.ID,
		}
		c.SendRoomUser(broadcast)
	} else {
		roomConfig.CurrWord = ""
		result := ws.Result{
			CMD:  "10011",
			Data: roomConfig,
			Msg:  "",
		}
		broadcast := ws.RoomBroadcast{
			RoomID:   roomId,
			Data:     result,
			ClientID: c.Client.ID,
		}
		c.SendRoomUser(broadcast)
	}

}
