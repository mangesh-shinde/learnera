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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var mongoClient *mongo.Client = db.ConnectMongoDB()

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

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>Welcome to LearnEra</h1>"))
}

func GetCourseById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	courseId := vars["courseid"]
	objectId, err := primitive.ObjectIDFromHex(courseId)
	utils.CheckError(err)

	coll := db.GetCollection(mongoClient, "courses")
	filter := bson.M{"_id": objectId}
	result := coll.FindOne(context.Background(), filter)
	var course models.Course
	err = result.Decode(&course)
	utils.CheckError(err)

	json.NewEncoder(w).Encode(course)
}

func GetAllCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	coll := db.GetCollection(mongoClient, "courses")
	filter := bson.D{{}}
	cur, err := coll.Find(context.Background(), filter)
	utils.CheckError(err)

	defer cur.Close(context.Background())

	var courses []models.Course
	//iterate over result returned by cursor
	for cur.Next(context.Background()) {
		var courseDoc models.Course
		err := cur.Decode(&courseDoc)
		utils.CheckError(err)
		courses = append(courses, courseDoc)
	}

	if len(courses) > 0 {
		json.NewEncoder(w).Encode(courses)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "No courses are available"})

}

func AddCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var course models.Course
	err := json.Unmarshal(body, &course)
	utils.CheckError(err)

	coll := db.GetCollection(mongoClient, "courses")

	result, err := coll.InsertOne(context.Background(), course)
	utils.CheckError(err)

	message := fmt.Sprintf("Inserted document successfully: ID: %s!", result.InsertedID)

	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func DeleteAllCourses(w http.ResponseWriter, r *http.Request) {
	coll := db.GetCollection(mongoClient, "courses")
	filter := bson.D{{}}
	deletedResult, err := coll.DeleteMany(context.Background(), filter)
	utils.CheckError(err)
	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("Deleted %d records from database!", deletedResult.DeletedCount)})
}

func DeleteCourseById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	courseId := vars["courseid"]
	objectId, err := primitive.ObjectIDFromHex(courseId)
	utils.CheckError(err)

	coll := db.GetCollection(mongoClient, "courses")
	filter := bson.M{"_id": objectId}
	deleteResult, err := coll.DeleteOne(context.Background(), filter)
	utils.CheckError(err)

	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("%d course deleted!", deleteResult.DeletedCount)})

}

func AddCourseToCart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	courseId := vars["courseid"]

	//Extract username from Request body
	bytes, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	utils.CheckError(err)

	var userDetails models.UserCredentials
	err = json.Unmarshal(bytes, &userDetails)
	utils.CheckError(err)

	username := userDetails.Username
	userFilter := bson.M{"username": username}
	userColl := db.GetCollection(mongoClient, "users")
	userResult := userColl.FindOne(context.Background(), userFilter)
	var user models.UserCredentials
	err = userResult.Decode(&user)
	if err == mongo.ErrNoDocuments {
		json.NewEncoder(w).Encode(map[string]string{"error": "User doesn't exist! Please verify."})
		return
	}
	utils.CheckError(err)

	//Check if courseId exists in db
	//If yes, then add it to user cart, else  throw error
	courseObjectId, err := primitive.ObjectIDFromHex(courseId)
	utils.CheckError(err)

	var course models.Course

	coll := db.GetCollection(mongoClient, "courses")
	filter := bson.M{"_id": courseObjectId}
	result := coll.FindOne(context.Background(), filter)
	err = result.Decode(&course)
	if err == mongo.ErrNoDocuments {
		json.NewEncoder(w).Encode(map[string]string{"error": "Course doesn't exist! Please verify."})
		return
	}
	utils.CheckError(err)

	profileColl := db.GetCollection(mongoClient, "profiles")
	updResult, err := profileColl.UpdateOne(context.TODO(), bson.M{"username": username}, bson.M{"$push": bson.M{"cart": courseObjectId}})
	utils.CheckError(err)

	json.NewEncoder(w).Encode(map[string]string{"message": fmt.Sprintf("%d Course added to Cart!", updResult.MatchedCount)})

}
