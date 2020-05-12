package handler

import (
	"encoding/json"
	"github/com/yuuki80code/game-server/cmd_router/handler/model"
	"github/com/yuuki80code/game-server/config/define"
	"github/com/yuuki80code/game-server/controller/response"
	"github/com/yuuki80code/game-server/redis"
	"github/com/yuuki80code/game-server/ws"
	"strings"
)

func Draw(c *ws.Context) {
	roomBoradcast := ws.RoomBroadcast{
		RoomID: c.Client.RoomID,
		Data: ws.Result{
			CMD:  c.CMD,
			Data: c.Data,
			Msg:  "",
		},
		ClientID: c.Client.ID,
	}
	c.Client.RoomBroadcast(roomBoradcast)
}

func Answer(c *ws.Context) {
	cache := redis.Cache
	var roomConfig define.RoomConfig
	var danmu model.Danmu
	roomId := c.Client.RoomID
	if roomId == "" {

	}
	roomConfigStr, _ := cache.HGet(define.RedisRoomConfigKey, roomId).Result()
	json.Unmarshal([]byte(roomConfigStr), &roomConfig)

	_ = c.Bind(&danmu)
	//如果弹幕是答案的话
	if danmu.Danmu == roomConfig.CurrWord {
		//自己不能回答自己！
		if c.Client.ID == roomConfig.CurrUser {
			danmu.IsAnswer = false
		} else {
			danmu.IsAnswer = true

			if roomConfig.Round == roomConfig.AllRound*2 {
				//最后一轮
				defer doGameOver(c, roomConfig)
			} else {
				//进行下一轮
				defer doNext(c, roomConfig)
			}

		}

	} else {
		danmu.IsAnswer = false
	}

	danmu.Danmu = strings.ReplaceAll(danmu.Danmu, roomConfig.CurrWord, "***")

	roomBroadcast := ws.RoomBroadcast{
		RoomID: c.Client.RoomID,
		Data: ws.Result{
			CMD:  c.CMD,
			Data: danmu,
			Msg:  "",
		},
		ClientID: c.Client.ID,
	}

	c.Client.RoomBroadcastAll(roomBroadcast)

}

func doNext(c *ws.Context, roomConfig define.RoomConfig) {
	cache := redis.Cache
	//告诉下一个用户
	roomConfig.Round += 1
	userStr, _ := cache.LIndex(c.Client.RoomID, int64((roomConfig.Round-1)%roomConfig.AllRound)).Result()
	var user response.UserInfoResp
	json.Unmarshal([]byte(userStr), &user)
	roomConfig.CurrUser = user.Account
	roomConfig.CurrWord = roomConfig.AllWord[roomConfig.Round-1].Word
	cache.HSet(define.RedisRoomConfigKey, c.Client.RoomID, roomConfig)
	toNextUser := ws.RoomBroadcast{
		RoomID: c.Client.RoomID,
		Data: ws.Result{
			CMD:  "11000",
			Data: model.CurrPlayer{CurrWord: roomConfig.CurrWord},
			Msg:  "",
		},
		ClientID: user.Account,
	}
	c.Client.RoomSendUser(toNextUser)

	//通知其他用户
	result := ws.Result{
		CMD: "10011",
		Data: model.OtherPlayer{
			CurrUser: roomConfig.CurrUser,
			Round:    roomConfig.Round,
		},
		Msg: "",
	}
	broadcast := ws.RoomBroadcast{
		RoomID:   c.Client.RoomID,
		Data:     result,
		ClientID: user.Account,
	}
	c.Client.RoomBroadcast(broadcast)

}

func doGameOver(c *ws.Context, roomConfig define.RoomConfig) {
	cache := redis.Cache
	cache.HDel(define.RedisRoomConfigKey, c.Client.RoomID)
	cache.Del(c.Client.RoomID)
	//通知所有玩家游戏结束
	result := ws.Result{
		CMD:  "13000",
		Data: "{}",
		Msg:  "",
	}
	broadcast := ws.RoomBroadcast{
		RoomID: c.Client.RoomID,
		Data:   result,
	}
	c.Client.RoomBroadcastAll(broadcast)
}
