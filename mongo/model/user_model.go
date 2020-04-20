package model

import (
	"context"
	"github/com/yuuki80code/game-server/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type UserModel struct {
	Account    string
	Username   string
	Password   string
	Avatar     string
	CreateTime time.Time
}

func (m *UserModel) getCollection() *mongo2.Collection {
	return mongo.MongoDatabase.Collection("user")
}

func (m *UserModel) InsertUser() error {
	result, err := m.getCollection().InsertOne(context.TODO(), m)
	if err != nil {
		log.Print(err)
	}
	log.Print(result)
	return err
}

func (m *UserModel) FindUserByAccount(account string) error {
	filter := bson.D{{"account", account}}
	err := m.getCollection().FindOne(context.TODO(), filter).Decode(m)
	return err
}
