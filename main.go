package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github/com/yuuki80code/game-server/cmd_router"
	"github/com/yuuki80code/game-server/config"
	"github/com/yuuki80code/game-server/mongo"
	"github/com/yuuki80code/game-server/router"
	"github/com/yuuki80code/game-server/ws"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	//初始化配置
	config.Load()
	//初始化mongo
	mongo.InitMongo()

	hub := ws.NewHub()
	go hub.Run()
	//初始化cmd
	cmd_router.InitRouter()
	r := gin.New()
	r.Use(gin.Recovery())
	//初始化路由
	router.InitRouter(r)
	r.Any("/ws", func(context *gin.Context) {
		ws.ServeWs(hub, context.Writer, context.Request)
	})
	//http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	//	ws.ServeWs(hub, w, r)
	//})
	r.Run(*addr)
	//err := http.ListenAndServe(*addr, nil)
	//if err != nil {
	//	log.Fatal("ListenAndServe: ", err)
	//}
}
