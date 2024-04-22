package repository

import (
	"context"
	"time"

	"dev.chaiyapluek.cloud.final.backend/src/config"
	"dev.chaiyapluek.cloud.final.backend/src/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmailRepository interface {
	Save(email *entity.Email) error
	GetNumberOfEmailSendWithInADay(to string) (int, error)
}

type emailRepository struct {
	collection *mongo.Collection
	cfg        *config.Collection
}

func NewEmailRepository(conn *mongo.Client, dbName string, cfg *config.Collection) EmailRepository {
	collection := conn.Database(dbName).Collection(cfg.Email)
	return &emailRepository{
		collection: collection,
		cfg:        cfg,
	}
}

func (r *emailRepository) Save(email *entity.Email) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	email.Id = nil
	_, err := r.collection.InsertOne(ctx, email)
	if err != nil {
		return err
	}

	return nil
}

func (r *emailRepository) GetNumberOfEmailSendWithInADay(to string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"to": to,
		"send_at": bson.M{
			"$gte": time.Now().Add(-24 * time.Hour),
		},
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
