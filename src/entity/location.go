package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Location struct {
	Id   *primitive.ObjectID `bson:"_id,omitempty"`
	Name string              `bson:"name"`
}

type Menu struct {
	Id             *primitive.ObjectID `bson:"_id,omitempty"`
	LocationId     *primitive.ObjectID `bson:"location_id"`
	Name           string              `bson:"name"`
	Description    string              `bson:"description"`
	Price          int                 `bson:"price"`
	IconImage      string              `bson:"iconImage"`
	ThumbnailImage string              `bson:"thumbnailImage"`
	Steps          []*Step             `bson:"steps"`
}

type Step struct {
	Name        string    `bson:"name"`
	Description string    `bson:"description"`
	Type        string    `bson:"type"`
	Required    bool      `bson:"required"`
	Min         int       `bson:"min"`
	Max         int       `bson:"max"`
	Options     []*Option `bson:"options"`
}

type Option struct {
	Name  string `db:"name"`
	Value string `db:"value"`
	Price int    `db:"price"`
}
