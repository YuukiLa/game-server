package controller

import (
	"github.com/gin-gonic/gin"
	"github/com/yuuki80code/game-server/controller/request"
	"github/com/yuuki80code/game-server/service"
	"github/com/yuuki80code/game-server/util"
)

type UserController struct {
	UserService *service.UserService
}

func (this *UserController) Register(c *gin.Context) {
	var user request.UserRequest
	err := c.ShouldBindJSON(&user)
	if err != nil {
		util.SendSimpleFailResp(c, 400, "参数错误")
		return
	}

	err = this.UserService.Register(user.Account, user.Password)
	if err != nil {
		util.SendSimpleFailResp(c, 500, err.Error())
		return
	}
	util.SendSimpleSuccessResp(c, "注册成功")
}

func (this *UserController) Login(c *gin.Context) {
	var user request.UserRequest
	err := c.ShouldBindJSON(&user)
	if err != nil {
		util.SendSimpleFailResp(c, 400, "参数错误")
		return
	}
	resp, err := this.UserService.Login(user.Account, user.Password)
	if err != nil {
		util.SendSimpleFailResp(c, 500, err.Error())
		return
	}
	util.SendDataSuccessResp(c, resp)
}
