package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type TodoModel struct {
	Id    primitive.ObjectID `bson:"_id" json:"id"`
	Title string             `bson:"title" json:"title" validate:"required"`
}
