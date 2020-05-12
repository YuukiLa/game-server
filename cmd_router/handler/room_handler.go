package handler

import (
	"encoding/json"
	"github/com/yuuki80code/game-server/cmd_router/handler/model"
	"github/com/yuuki80code/game-server/config/define"
	"github/com/yuuki80code/game-server/controller/response"
	model2 "github/com/yuuki80code/game-server/mongo/model"
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

		len, _ := cache.LLen(roomId).Result()

		roomConfig.Status = 1
		roomConfig.CurrUser = c.Client.ID

		roomConfig.Round = 1
		roomConfig.AllRound = int(len)

		words, err := new(model2.WordModel).GetRandomWord(len * 2)
		if err != nil {
			log.Println("获取词失败")
		}
		roomConfig.AllWord = words
		roomConfig.CurrWord = words[roomConfig.Round-1].Word
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
			Data: model.CurrPlayer{CurrWord: roomConfig.CurrWord},
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
			CMD: "10011",
			Data: model.OtherPlayer{
				CurrUser: roomConfig.CurrUser,
				Round:    roomConfig.Round,
			},
			Msg: "",
		}
		broadcast := ws.RoomBroadcast{
			RoomID:   roomId,
			Data:     result,
			ClientID: c.Client.ID,
		}
		c.SendRoomUser(broadcast)
	}

}

func ExitRoom(c *ws.Context) {
	cache := redis.Cache
	var roomConfig define.RoomConfig

	roomId := c.Client.RoomID
	if roomId == "" {

	}
	roomConfigStr, _ := cache.HGet(define.RedisRoomConfigKey, roomId).Result()
	json.Unmarshal([]byte(roomConfigStr), &roomConfig)
	//解散房间
	if roomConfig.Master == c.Client.ID {
		cache.HDel(define.RedisRoomConfigKey, roomId)
		cache.Del(roomId)
		result := ws.Result{
			CMD:  "10008",
			Data: "{}",
			Msg:  "",
		}
		broadcast := ws.RoomBroadcast{
			RoomID:   roomId,
			Data:     result,
			ClientID: "",
		}
		c.SendRoomBroadcastAll(broadcast)
	} else {
		users := make([]response.UserInfoResp, 0)
		data, err := cache.LRange(roomId, 0, -1).Result()
		if err != nil {
			log.Println(err)
			c.SendString(err.Error())
			return
		}
		cache.Del(roomId)
		for _, value := range data {
			var user response.UserInfoResp
			json.Unmarshal([]byte(value), &user)
			if user.Account != c.Client.ID {
				users = append(users, user)
				cache.RPush(roomId, user)
			}
		}

		result := ws.Result{
			CMD:  "10001",
			Data: users,
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
