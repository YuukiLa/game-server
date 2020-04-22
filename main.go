package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github/com/yuuki80code/game-server/cmd_router"
	"github/com/yuuki80code/game-server/config"
	"github/com/yuuki80code/game-server/mongo"
	"github/com/yuuki80code/game-server/redis"
	"github/com/yuuki80code/game-server/router"
	"github/com/yuuki80code/game-server/ws"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	//初始化配置
	config.Load()
	//初始化mongo
	mongo.InitMongo()
	//初始化redis
	redis.InitRedis()
	err := redis.Cache.Set("test", "test", 100)
	if err != nil {
		log.Print(err)
	}
	hub := ws.NewHub()
	go hub.Run()
	//初始化cmd
	cmd_router.InitRouter()
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(Cors())
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

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
