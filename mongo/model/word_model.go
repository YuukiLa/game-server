package model

import (
	"context"
	"encoding/json"
	"github/com/yuuki80code/game-server/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"log"
)

type WordModel struct {
	Word  string
	Hint1 string
	Hint2 string
}

func (s WordModel) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s WordModel) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

func (m *WordModel) getCollection() *mongo2.Collection {
	return mongo.MongoDatabase.Collection("word")
}

func (m *WordModel) GetRandomWord(size int64) (words []WordModel, err error) {
	pipeline := []bson.D{bson.D{{"$sample", bson.D{{"size", size}}}}}
	cursor, err := m.getCollection().Aggregate(context.Background(), pipeline)
	if err != nil {
		log.Println(err)
	}
	err = cursor.All(context.Background(), &words)
	//for cursor.TryNext(context.Background()) {
	//
	//}
	if err != nil {
		log.Println(err)
	}
	return
}
