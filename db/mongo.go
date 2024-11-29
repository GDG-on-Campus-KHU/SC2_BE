package db

import "go.mongodb.org/mongo-driver/mongo"

var MongoClient *mongo.Client

func SetMongoClient(client *mongo.Client) {
	MongoClient = client
}

func GetMongoClient() *mongo.Client {
	return MongoClient
}
