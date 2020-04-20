package middleware

import (
	"github.com/gin-gonic/gin"
	"github/com/yuuki80code/game-server/util"
	"net/http"
)

func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		checkToken(c)
	}
}

func checkToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.String(http.StatusUnauthorized, "")
		c.Abort()
		return
	}
	if token == "undefined" {
		c.String(http.StatusUnauthorized, "")
		c.Abort()
		return
	}
	account, err := util.Dncrypt(token)
	if err != nil {
		c.String(http.StatusUnauthorized, "token过期")
		c.Abort()
		return
	}
	c.Set("account", account)
	c.Set("auth_token", token)
	c.Next()
}
