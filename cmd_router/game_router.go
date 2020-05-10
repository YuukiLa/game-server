package cmd_router

import (
	"github/com/yuuki80code/game-server/cmd_router/handler"
	"github/com/yuuki80code/game-server/ws"
)

func initGameRouter() {
	ws.AddHandler("10011", handler.GameParam)
	ws.AddHandler("10012", handler.Draw)
	ws.AddHandler("10013", handler.Answer)
}
