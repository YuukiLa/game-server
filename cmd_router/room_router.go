package cmd_router

import (
	"github/com/yuuki80code/game-server/cmd_router/handler"
	"github/com/yuuki80code/game-server/ws"
)

func initRoomRouter() {
	ws.AddHandler("10001", handler.EnterRoom)
}
