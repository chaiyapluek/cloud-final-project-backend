package repository

import (
	"context"
	"log"
	"time"

	"dev.chaiyapluek.cloud.final.backend/src/config"
	"dev.chaiyapluek.cloud.final.backend/src/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepository interface {
	GetLoginAttemptById(id *primitive.ObjectID) (*entity.LoginAttempt, error)
	GetRegisterAttemptById(id *primitive.ObjectID) (*entity.RegisterAttempt, error)
	CreateLoginAttempt(e *entity.LoginAttempt) error
	CreateRegisterAttempt(e *entity.RegisterAttempt) error
}

type authRepository struct {
	collection *mongo.Collection
	cfg        *config.Collection
}

func NewAuthRepository(db *mongo.Client, dbName string, cfg *config.Collection) AuthRepository {
	collection := db.Database(dbName).Collection(cfg.AuthAttempt)
	return &authRepository{
		collection: collection,
		cfg:        cfg,
	}
}

func (r *authRepository) GetLoginAttemptById(id *primitive.ObjectID) (*entity.LoginAttempt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var loginAttempt entity.LoginAttempt
	err := r.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&loginAttempt)

	if err != nil {
		log.Println("repository, get login attempt", err)
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &loginAttempt, nil
}

func (r *authRepository) CreateLoginAttempt(e *entity.LoginAttempt) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	e.Id = nil
	result, err := r.collection.InsertOne(ctx, e)
	if err != nil {
		log.Println("repository, create login attempt", err)
		return err
	}

	log.Println("repository, create login attempt success", result.InsertedID)
	oid := result.InsertedID.(primitive.ObjectID)
	e.Id = &oid

	return nil
}

func (r *authRepository) GetRegisterAttemptById(id *primitive.ObjectID) (*entity.RegisterAttempt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var registerAttempt entity.RegisterAttempt
	err := r.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&registerAttempt)

	if err != nil {
		log.Println("repository, get register attempt", err)
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &registerAttempt, nil
}

func (r *authRepository) CreateRegisterAttempt(e *entity.RegisterAttempt) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	e.Id = nil
	result, err := r.collection.InsertOne(ctx, e)
	if err != nil {
		log.Println("repository, create register attempt", err)
		return err
	}

	log.Println("repository, create register attempt success", result.InsertedID)
	oid := result.InsertedID.(primitive.ObjectID)
	e.Id = &oid

	return nil
}
