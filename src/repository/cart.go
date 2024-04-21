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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CartRepository interface {
	GetByCartId(cartId *primitive.ObjectID) (*entity.Cart, error)
	GetCartByUserId(userId, locationId *primitive.ObjectID) (*entity.Cart, error)
	CreateCart(e *entity.Cart) error
	AddCartItem(cartId *primitive.ObjectID, cartItem *entity.CartItem) (*entity.Cart, error)
	DeleteCartItem(cartId *primitive.ObjectID, itemId int) (*entity.Cart, error)
	DeleteCartById(cartId *primitive.ObjectID) error
}

type cartRepositoryImpl struct {
	collection *mongo.Collection
	cfg        *config.Collection
}

func NewCartRepository(db *mongo.Client, dbName string, cfg *config.Collection) CartRepository {
	collection := db.Database(dbName).Collection(cfg.Cart)
	return &cartRepositoryImpl{
		collection: collection,
		cfg:        cfg,
	}
}

func (r *cartRepositoryImpl) GetByCartId(cartId *primitive.ObjectID) (*entity.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cart entity.Cart
	err := r.collection.FindOne(ctx, bson.M{
		"_id": cartId,
	}).Decode(&cart)

	if err != nil {
		log.Println("repository, get cart by id", err)
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepositoryImpl) GetCartByUserId(userId, locationId *primitive.ObjectID) (*entity.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"user_id":     userId,
				"location_id": locationId,
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$items",
				"preserveNullAndEmptyArrays": true,
			},
		},
		{
			"$lookup": bson.M{
				"from":         r.cfg.Menu,
				"localField":   "items.menu_id",
				"foreignField": "_id",
				"pipeline": []bson.M{
					{
						"$project": bson.M{
							"name": 1,
						},
					},
				},
				"as": "items.menu_name",
			},
		},
		{
			"$set": bson.M{
				"items.menu_name": bson.M{
					"$arrayElemAt": []interface{}{"$items.menu_name.name", 0},
				},
			},
		},
		{
			"$group": bson.M{
				"_id": "$_id",
				"user_id": bson.M{
					"$first": "$user_id",
				},
				"location_id": bson.M{
					"$first": "$location_id",
				},
				"session_id": bson.M{
					"$first": "$session_id",
				},
				"items": bson.M{
					"$push": "$items",
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         r.cfg.Location,
				"localField":   "location_id",
				"foreignField": "_id",
				"as":           "location_name",
			},
		},
		{
			"$set": bson.M{
				"location_name": bson.M{
					"$arrayElemAt": []interface{}{"$location_name.name", 0},
				},
			},
		},
	}

	var carts []*entity.Cart

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("repository, get cart by user id", err)
		return nil, err
	}

	if err = cursor.All(ctx, &carts); err != nil {
		log.Println("repository, get cart by user id", err)
		return nil, err
	}

	if len(carts) == 0 {
		return nil, nil
	}

	return carts[0], nil
}

func (r *cartRepositoryImpl) CreateCart(e *entity.Cart) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	e.Id = nil
	result, err := r.collection.InsertOne(ctx, e)
	if err != nil {
		log.Println("repository, create cart", err)
		return err
	}

	oid := result.InsertedID.(primitive.ObjectID)
	e.Id = &oid
	return nil
}

func (r *cartRepositoryImpl) AddCartItem(cartId *primitive.ObjectID, cartItem *entity.CartItem) (*entity.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	after := options.After
	ops := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	var cart entity.Cart
	err := r.collection.FindOneAndUpdate(ctx, bson.M{
		"_id": cartId,
	}, bson.M{
		"$push": bson.M{
			"items": cartItem,
		},
	}, &ops).Decode(&cart)
	if err != nil {
		log.Println("repository, add cart item", err)
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepositoryImpl) DeleteCartItem(cartId *primitive.ObjectID, itemId int) (*entity.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	after := options.After
	ops := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	var cart entity.Cart
	err := r.collection.FindOneAndUpdate(ctx, bson.M{
		"_id": cartId,
	}, bson.M{
		"$pull": bson.M{
			"items": bson.M{
				"item_id": itemId,
			},
		},
	}, &ops).Decode(&cart)

	if err != nil {
		log.Println("repository, delete cart item", err)
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepositoryImpl) DeleteCartById(cartId *primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{
		"_id": cartId,
	})
	if err != nil {
		log.Println("repository, delete cart by id", err)
		return err
	}

	return nil
}
