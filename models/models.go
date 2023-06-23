package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Course struct {
	Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CourseName string             `json:"coursename"`
	Author     string             `json:"author"`
	Price      int                `json:"price"`
}
