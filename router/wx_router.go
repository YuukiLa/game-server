package router

import "github/com/yuuki80code/game-server/controller"

var WxController *controller.WxController

func initWxRouter() {
	apiRouter.POST("/code",WxController.Code)
}
