package handler

import (
	"encoding/json"
	"github/com/yuuki80code/game-server/cmd_router/handler/model"
	"github/com/yuuki80code/game-server/config/define"
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
		}

	} else {
		danmu.IsAnswer = false
	}

	danmu.Danmu = strings.ReplaceAll(danmu.Danmu, roomConfig.CurrWord, "***")

	roomBoradcast := ws.RoomBroadcast{
		RoomID: c.Client.RoomID,
		Data: ws.Result{
			CMD:  c.CMD,
			Data: danmu,
			Msg:  "",
		},
		ClientID: c.Client.ID,
	}

	c.Client.RoomBroadcastAll(roomBoradcast)

}
