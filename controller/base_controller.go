package controller

import "github.com/gin-gonic/gin"

func GetUserAccount(c *gin.Context) string {
	account, _ := c.Get("account")
	return account.(string)
}
