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

func (this *UserController) Userinfo(c *gin.Context) {
	account := GetUserAccount(c)
	resp, err := this.UserService.UserInfo(account)
	if err != nil {
		util.SendSimpleFailResp(c, 401, "用户不存在")
		return
	}
	util.SendDataSuccessResp(c, resp)
}

func (this *UserController) UpdateUserName(c *gin.Context) {
	account := GetUserAccount(c)
	username := c.PostForm("username")
	err := this.UserService.UpdateUserName(account, username)
	if err != nil {
		util.SendSimpleFailResp(c, 400, err.Error())
		return
	}
	util.SendSimpleSuccessResp(c, "修改成功")
}

func (this *UserController) UpdateUserAvatar(c *gin.Context) {
	account := GetUserAccount(c)
	avatar := c.PostForm("avatar")
	err := this.UserService.UpdateAvatar(account, avatar)
	if err != nil {
		util.SendSimpleFailResp(c, 400, err.Error())
		return
	}
	util.SendSimpleSuccessResp(c, "修改成功")
}
