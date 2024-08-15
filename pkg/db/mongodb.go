package db

import (
	"api-gateway/pkg/logger"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Client

func ConnectMongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://admin:admin123@localhost:27017")

	var err error
	MongoDB, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logger.Fatal("無法連接到MongoDB: ", err)
	}

	err = MongoDB.Ping(context.TODO(), nil)
	if err != nil {
		logger.Fatal("無法連接到MongoDB: ", err)
	}

}

func CloseMongoDB() {
	defer MongoDB.Disconnect(context.TODO())
}
