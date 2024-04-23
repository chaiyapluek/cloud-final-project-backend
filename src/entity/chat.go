package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AIAnswer struct {
	Explanation string `json:"explanation"`
	Bread       string `json:"bread"`
	BreadSize   string `json:"bread_size"`
	Menu        string `json:"menu"`
	Vegetables  string `json:"vegetables"`
	Sauces      string `json:"sauces"`
	AddOns      string `json:"add_ons"`
	Meal        string `json:"meal"`
	TotalPrice  string `json:"total_price"`
}

type Chat struct {
	Id         *primitive.ObjectID `bson:"_id,omitempty"`
	LocationId *primitive.ObjectID `bson:"location_id"`
	UserId     *primitive.ObjectID `bson:"user_id"`
	Type       int                 `bson:"type"`
	Sender     int                 `bson:"sender"`
	RawContent string              `bson:"raw_content"`
	Content    string              `bson:"content"`
	SentAt     time.Time           `bson:"sent_at"`
}
