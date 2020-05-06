package redis

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github/com/yuuki80code/game-server/config"
)

var Cache *redis.Client

func InitRedis() {

	Cache = redis.NewClient(&redis.Options{
		Addr:     config.Configer.Redis.Url,
		Password: config.Configer.Redis.Password,
		DB:       config.Configer.Redis.Db,
	})

	pong, err := Cache.Ping().Result()
	fmt.Println(pong, err)
}
