// db/db.go
package db

import (
    "context"
    "fmt"
    "os"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
)

func ConnectDB() (*mongo.Client, error) {
    userURI := os.Getenv("MONGO_URI")
    serverAPI := options.ServerAPI(options.ServerAPIVersion1)
    opts := options.Client().ApplyURI(userURI).SetServerAPIOptions(serverAPI)

    client, err := mongo.Connect(context.TODO(), opts)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
    }

    if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
        return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
    }

    fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
    return client, nil
}

func DisconnectDB(client *mongo.Client) error {
    err := client.Disconnect(context.TODO())
    if err != nil {
        return fmt.Errorf("failed to disconnect from MongoDB: %v", err)
    }

    fmt.Println("Successfully disconnected from MongoDB!")
    return nil
}