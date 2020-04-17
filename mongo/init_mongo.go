package mongo

import (
	"context"
	"fmt"
	"github/com/yuuki80code/game-server/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)


var (
	MongoClient *mongo.Client
	MongoDatabase *mongo.Database
)

func InitMongo() {
	linkUri := fmt.Sprintf("mongodb://%s:%s@%s:%s",config.Configer.Mongo.User,config.Configer.Mongo.Password,config.Configer.Mongo.Url,config.Configer.Mongo.Port)
	clientOptions := options.Client().ApplyURI(linkUri)
	MongoClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// 检查连接
	err = MongoClient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	MongoDatabase = MongoClient.Database(config.Configer.Mongo.Database)
}
