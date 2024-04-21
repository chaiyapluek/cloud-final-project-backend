package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       *primitive.ObjectID `bson:"_id,omitempty"`
	Email    string              `bson:"email"`
	Password string              `bson:"password"`
	Name     string              `bson:"name"`
}
