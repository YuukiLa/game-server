package router

import "github.com/gin-gonic/gin"

var (
	apiRouter *gin.RouterGroup
)

func InitRouter(r *gin.Engine) {
	apiRouter = r.Group("/api")
	initWxRouter()
}
