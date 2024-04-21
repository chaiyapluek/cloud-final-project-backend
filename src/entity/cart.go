package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type ItemStep struct {
	Step    string   `bson:"step"`
	Options []string `bson:"options"`
}

type CartItem struct {
	MenuId     *primitive.ObjectID `bson:"menu_id"`
	MenuName   string              `bson:"menu_name"`
	ItemId     int                 `bson:"item_id"`
	Quantity   int                 `bson:"quantity"`
	TotalPrice int                 `bson:"total_price"`
	Steps      []ItemStep          `bson:"steps"`
}

type Cart struct {
	Id           *primitive.ObjectID `bson:"_id,omitempty"`
	UserId       *primitive.ObjectID `bson:"user_id"`
	LocationId   *primitive.ObjectID `bson:"location_id"`
	LocationName string              `bson:"location_name"`
	SessionId    *string             `bson:"session_id"`
	Items        []*CartItem
}
