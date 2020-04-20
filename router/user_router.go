package router

import "github/com/yuuki80code/game-server/controller"

var UserController = &controller.UserController{}

func initUserRouter() {
	apiRouter.POST("/register", UserController.Register)
	apiRouter.POST("/login", UserController.Login)
}
