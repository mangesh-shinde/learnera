package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mangesh-shinde/learnera/db"
	"github.com/mangesh-shinde/learnera/models"
	"github.com/mangesh-shinde/learnera/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func AddUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	byteData, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	utils.CheckError(err)

	var profile models.UserProfile
	err = json.Unmarshal(byteData, &profile)
	utils.CheckError(err)

	coll := db.GetCollection(mongoClient, "profiles")
	insertResult, err := coll.InsertOne(context.Background(), profile)
	utils.CheckError(err)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Profile added to DB with id: %s", insertResult.InsertedID)})

}

func GetAllUserProfiles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//Get user profile from mongodb
	var profiles []models.UserProfile
	coll := db.GetCollection(mongoClient, "profiles")
	filter := bson.M{}
	cur, err := coll.Find(context.Background(), filter)
	utils.CheckError(err)

	for cur.Next(context.Background()) {
		var profile models.UserProfile
		err := cur.Decode(&profile)
		utils.CheckError(err)
		profiles = append(profiles, profile)
	}

	if len(profiles) <= 0 {
		w.WriteHeader(http.StatusNoContent)
		json.NewEncoder(w).Encode(map[string]string{"error": "No user profiles are available!"})
		return
	}

	json.NewEncoder(w).Encode(profiles)

}

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	username := vars["username"]

	//Get user profile from mongodb
	var profile models.UserProfile
	coll := db.GetCollection(mongoClient, "profiles")
	filter := bson.M{"username": username}
	err := coll.FindOne(context.Background(), filter).Decode(&profile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusNoContent)
			json.NewEncoder(w).Encode(map[string]string{"error": "User profile is unavailable!"})
			return
		}
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(profile)

}

func AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	byteData, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	utils.CheckError(err)

	var userCredentials models.UserCredentials
	err = json.Unmarshal(byteData, &userCredentials)
	utils.CheckError(err)

	// encrypt password and save
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userCredentials.Password), bcrypt.DefaultCost)
	utils.CheckError(err)

	userCredentials.Password = string(hashedPassword)

	coll := db.GetCollection(mongoClient, "users")
	_, err = coll.InsertOne(context.Background(), userCredentials)
	utils.CheckError(err)

	json.NewEncoder(w).Encode(map[string]string{"message": "User added successfully!"})

}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//grab username and password from request body
	byteData, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	utils.CheckError(err)

	var userCredentials models.UserCredentials
	err = json.Unmarshal(byteData, &userCredentials)
	utils.CheckError(err)

	// Check if credentials match in mongodb
	var result models.UserCredentials
	coll := db.GetCollection(mongoClient, "users")
	filter := bson.M{"username": userCredentials.Username}
	err = coll.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
			return
		}
		log.Fatal(err)
	}

	inputPassword := userCredentials.Password
	hashedPassword := result.Password

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Incorrect Password!"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful!"})
}
