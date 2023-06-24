package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Course struct {
	Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CourseName string             `json:"coursename"`
	Author     string             `json:"author"`
	Price      int                `json:"price"`
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserProfile struct {
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}
