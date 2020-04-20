package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func SendSimpleSuccessResp(c *gin.Context, msg string) {
	SendSuccessResp(c, msg, nil)
}

func SendDataSuccessResp(c *gin.Context, data interface{}) {
	SendSuccessResp(c, "success", data)
}

func SendSuccessResp(c *gin.Context, msg string, data interface{}) {
	resp := Resp{
		Code: 200,
		Msg:  msg,
		Data: data,
	}
	c.JSON(http.StatusOK, resp)
}

func SendSimpleFailResp(c *gin.Context, code int, msg string) {
	resp := Resp{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
	c.JSON(http.StatusOK, resp)
}
