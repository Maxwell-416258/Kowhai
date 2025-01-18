package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"kowhai/global"
	"log"
)

func InitMongo() *mongo.Client {
	host := global.Config.Mongo.Host
	port := global.Config.Mongo.Port
	user := global.Config.Mongo.User
	password := global.Config.Mongo.Password
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d", user, password, host, port)
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		global.Logger.Fatalf("Failed to connect to mongo: %v", err)
	}

	// 确保连接成功
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
		global.Logger.Fatalf("Failed to ping to mongo: %v", err)
	}

	fmt.Println("Connected to MongoDB!")

	return client

}
