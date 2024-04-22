package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Email struct {
	Id     *primitive.ObjectID `bson:"_id,omitempty"`
	To     string              `bson:"to"`
	SendAt time.Time           `bson:"send_at"`
}
