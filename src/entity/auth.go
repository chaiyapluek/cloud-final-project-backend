package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginAttempt struct {
	Id     *primitive.ObjectID `bson:"_id,omitempty"`
	Email  string              `bson:"email"`
	Code   string              `bson:"code"`
	Expire time.Time           `bson:"expire"`
}

type RegisterAttempt struct {
	Id       *primitive.ObjectID `bson:"_id,omitempty"`
	Email    string              `bson:"email"`
	Password string              `bson:"password"`
	Name     string              `bson:"name"`
	Code     string              `bson:"code"`
	Expire   time.Time           `bson:"expire"`
}
