package controller

import (
	"github.com/gin-gonic/gin"
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
	err = cache.HSetNX(roomId, account, info).Err()
	if err != nil {
		util.SendSimpleFailResp(c, 400, "创建房间失败")
		return
	}
	util.SendDataSuccessResp(c, roomId)
}

func (this *RoomController) EnterRoom(c *gin.Context) {
	cache := redis.Cache
	roomId := c.PostForm("roomId")
	num, err := cache.HLen(roomId).Result()
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
	_, err = cache.HSet(roomId, account, info).Result()
	if err != nil {
		util.SendSimpleFailResp(c, 400, "加入房间失败")
		return
	}
	util.SendDataSuccessResp(c, roomId)
}
