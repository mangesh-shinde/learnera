package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetMongoCreds() MongoCredentials {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoUsername := os.Getenv("MONGO_USERNAME")
	mongoPassword := os.Getenv("MONGO_PASSWORD")

	if mongoUsername == "" {
		log.Fatal("Mongo Username is not available. Please check .env file")
	}

	if mongoPassword == "" {
		log.Fatal("Mongo Password is not available. Please check .env file")
	}

	var MongoCreds MongoCredentials
	MongoCreds.MongoUsername = mongoUsername
	MongoCreds.MongoPassword = mongoPassword

	return MongoCreds
}

type MongoCredentials struct {
	MongoUsername string
	MongoPassword string
}
