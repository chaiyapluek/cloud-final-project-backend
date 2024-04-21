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

type UserRepository interface {
	GetById(id *primitive.ObjectID) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	CreateUser(user *entity.User) error
}

type userRepositoryImpl struct {
	collection *mongo.Collection
	cfg        *config.Collection
}

func NewUserRepository(db *mongo.Client, dbName string, cfg *config.Collection) UserRepository {
	collection := db.Database(dbName).Collection(cfg.User)
	return &userRepositoryImpl{
		collection: collection,
		cfg:        cfg,
	}
}

func (r *userRepositoryImpl) GetById(id *primitive.ObjectID) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{
		"_id": id,
	}).Decode(&user)

	if err != nil {
		log.Println("repository, get user by id", err)
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil

}

func (r *userRepositoryImpl) GetUserByEmail(email string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{
		"email": email,
	}).Decode(&user)

	if err != nil {
		log.Println("repository, get user by email", err)
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryImpl) CreateUser(user *entity.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user.Id = nil
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		log.Println("repository, create user", err)
		return err
	}

	log.Println("repository, create user", result.InsertedID)
	oid := result.InsertedID.(primitive.ObjectID)
	user.Id = &oid

	return nil
}
