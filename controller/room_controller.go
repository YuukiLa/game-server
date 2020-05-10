package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github/com/yuuki80code/game-server/config/define"
	"github/com/yuuki80code/game-server/redis"
	"github/com/yuuki80code/game-server/service"
	"github/com/yuuki80code/game-server/util"
	"log"
)

type RoomController struct {
	UserService *service.UserService
}

func (this *RoomController) CreateRoom(c *gin.Context) {
	cache := redis.Cache
	roomId := util.GetRandomNumber(6)
	account := GetUserAccount(c)
	info, err := this.UserService.UserInfo(account)
	if err != nil {
		util.SendSimpleFailResp(c, 400, "创建房间失败")
		return
	}
	//err = cache.RPush(roomId,info)
	err = cache.LPush(roomId, info).Err()
	if err != nil {
		util.SendSimpleFailResp(c, 400, "创建房间失败")
		return
	}
	roomConfig := define.RoomConfig{
		Master:   account,
		Status:   0,
		CurrWord: "",
		CurrUser: "",
		Round:    0,
	}

	_, err = cache.HSetNX(define.RedisRoomConfigKey, roomId, roomConfig).Result()
	if err != nil {
		log.Println(err)
	}
	util.SendDataSuccessResp(c, roomId)
}

func (this *RoomController) EnterRoom(c *gin.Context) {
	cache := redis.Cache
	roomId := c.PostForm("roomId")
	var roomConfig define.RoomConfig
	roomConfigStr, err := cache.HGet(define.RedisRoomConfigKey, roomId).Result()
	if err != nil {
		util.SendSimpleFailResp(c, 400, "找不到房间，请确认房间号正确")
		log.Println(err)
		return
	}
	json.Unmarshal([]byte(roomConfigStr), &roomConfig)
	if roomConfig.Status == 1 {
		util.SendSimpleFailResp(c, 400, "该房间已开始游戏无法加入")
		log.Println(err)
		return
	}
	num, err := cache.LLen(roomId).Result()
	if err != nil {
		util.SendSimpleFailResp(c, 400, "找不到房间，请确认房间号正确")
		log.Println(err)
		return
	}
	if num < 1 {
		util.SendSimpleFailResp(c, 400, "找不到房间，请确认房间号正确")
		return
	}
	if num > 5 {
		util.SendSimpleFailResp(c, 400, "房间已满")
		return
	}
	account := GetUserAccount(c)
	info, err := this.UserService.UserInfo(account)
	if err != nil {
		util.SendSimpleFailResp(c, 400, "加入房间失败")
		return
	}
	//err = cache.RPush(roomId,info)
	_, err = cache.LPush(roomId, info).Result()
	if err != nil {
		util.SendSimpleFailResp(c, 400, "加入房间失败")
		return
	}

	util.SendDataSuccessResp(c, roomId)
}
