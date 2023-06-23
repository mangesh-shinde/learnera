package db

import (
	"context"
	"fmt"
	"log"

	"github.com/mangesh-shinde/learnera/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB() *mongo.Client {

	mongocreds := utils.GetMongoCreds()

	ConnectionString := fmt.Sprintf("mongodb+srv://%s:%s@cluster0.x8ue6.mongodb.net/?retryWrites=true&w=majority", mongocreds.MongoUsername, mongocreds.MongoPassword)
	fmt.Println(ConnectionString)

	opts := options.Client().ApplyURI(ConnectionString)
	Client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	return Client

}

func GetCollection(cl *mongo.Client, collectionName string) *mongo.Collection {
	Collection := cl.Database("learnera").Collection(collectionName)
	return Collection
}

func DisconnectMongo(cl *mongo.Client) {
	err := cl.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Disconnected to MongoDB!")
}
