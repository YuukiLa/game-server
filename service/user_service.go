package service

import (
	"errors"
	"github/com/yuuki80code/game-server/controller/response"
	"github/com/yuuki80code/game-server/mongo/model"
	"github/com/yuuki80code/game-server/util"
	"time"
)

type UserService struct {
}

func (this *UserService) Register(account, password string) error {
	user := new(model.UserModel)
	err := user.FindUserByAccount(account)
	if err == nil {
		return errors.New("账号已被注册")
	}
	encodePwd := util.HashPwd(password)
	user = &model.UserModel{
		Account:    account,
		Username:   account,
		Password:   encodePwd,
		Avatar:     "1.webp",
		CreateTime: time.Now(),
	}
	err = user.InsertUser()
	return err
}

func (this *UserService) Login(account, password string) (*response.LoginResp, error) {
	user := new(model.UserModel)
	err := user.FindUserByAccount(account)
	if err != nil {
		return nil, errors.New("账号不存在")
	}
	if !util.EqualsPwd(password, user.Password) {
		return nil, errors.New("账号或密码错误")
	}
	token, _ := util.Encrypt(account)
	resp := response.LoginResp{
		Token: token,
		UserInfoResp: response.UserInfoResp{
			Account:  account,
			Username: user.Username,
			Avatar:   user.Avatar,
		},
	}
	return &resp, nil
}

func (this *UserService) UserInfo(account string) (*response.UserInfoResp, error) {
	user := new(model.UserModel)
	err := user.FindUserByAccount(account)
	if err != nil {
		return nil, errors.New("账号不存在")
	}
	userInfoResp := response.UserInfoResp{
		Username: user.Username,
		Avatar:   user.Avatar,
	}
	return &userInfoResp, nil
}

func (this *UserService) UpdateUserName(account string, username string) error {
	user := new(model.UserModel)
	return user.UpdateUserName(account, username)
}

func (this *UserService) UpdateAvatar(account string, username string) error {
	user := new(model.UserModel)
	return user.UpdateUserAvatar(account, username)
}
