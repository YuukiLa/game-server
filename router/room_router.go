package router

import (
	"github/com/yuuki80code/game-server/controller"
	"github/com/yuuki80code/game-server/middleware"
)

var RoomController = &controller.RoomController{}

func initRoomRouter() {
	apiRouter.POST("/room", middleware.CheckToken(), RoomController.CreateRoom)
	apiRouter.POST("/room/enter", middleware.CheckToken(), RoomController.EnterRoom)

}
