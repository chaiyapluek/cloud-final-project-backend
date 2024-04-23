package repository

import (
	"context"
	"time"

	"dev.chaiyapluek.cloud.final.backend/src/config"
	"dev.chaiyapluek.cloud.final.backend/src/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepository interface {
	GetUserChat(userId *primitive.ObjectID, locationId *primitive.ObjectID) ([]*entity.Chat, error)
	Insert(e *entity.Chat) error
	InsertMany(e []*entity.Chat) error
	DeleteChat(userId *primitive.ObjectID, locationId *primitive.ObjectID) error
}

type chatRepository struct {
	collection *mongo.Collection
	cfg        *config.Collection
}

func NewChatRepository(conn *mongo.Client, dbName string, cfg *config.Collection) ChatRepository {
	collection := conn.Database(dbName).Collection(cfg.Chat)
	return &chatRepository{
		collection: collection,
		cfg:        cfg,
	}
}

func (r *chatRepository) GetUserChat(userId *primitive.ObjectID, locationId *primitive.ObjectID) ([]*entity.Chat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ops := options.Find().SetSort(bson.D{{"sent_at", 1}})
	cursor, err := r.collection.Find(ctx, bson.M{
		"user_id":     userId,
		"location_id": locationId,
	}, ops)
	if err != nil {
		return nil, err
	}

	var chats []*entity.Chat
	if err = cursor.All(ctx, &chats); err != nil {
		return nil, err
	}

	return chats, nil
}

func (r *chatRepository) Insert(e *entity.Chat) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	e.Id = nil
	_, err := r.collection.InsertOne(ctx, e)
	if err != nil {
		return err
	}

	return nil
}

func (r *chatRepository) InsertMany(e []*entity.Chat) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newChats := []interface{}{}
	for _, v := range e {
		v.Id = nil
		newChats = append(newChats, v)
	}
	_, err := r.collection.InsertMany(ctx, newChats)
	if err != nil {
		return err
	}

	return nil
}

func (r *chatRepository) DeleteChat(userId *primitive.ObjectID, locationId *primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteMany(ctx, bson.M{
		"user_id":     userId,
		"location_id": locationId,
	})
	if err != nil {
		return err
	}

	return nil
}
