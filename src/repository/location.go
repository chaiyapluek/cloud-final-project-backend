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

type LocationRepository interface {
	GetAllLocation() ([]*entity.Location, error)
	GetLocationById(id *primitive.ObjectID) (*entity.Location, []*entity.Menu, error)
	GetMenuItmes(locationId, menuId *primitive.ObjectID) (*entity.Menu, error)
}

type locationRepositoryImpl struct {
	locationCollection *mongo.Collection
	menuCollection     *mongo.Collection
	cfg                *config.Collection
}

func NewLocationRepository(db *mongo.Client, dbName string, cfg *config.Collection) LocationRepository {
	locationCollection := db.Database(dbName).Collection(cfg.Location)
	menuCollection := db.Database(dbName).Collection(cfg.Menu)
	return &locationRepositoryImpl{
		locationCollection: locationCollection,
		menuCollection:     menuCollection,
		cfg:                cfg,
	}
}

func (r *locationRepositoryImpl) GetAllLocation() ([]*entity.Location, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.locationCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("repository, get all location", err)
		return nil, err
	}

	var result []*entity.Location
	if err = cursor.All(ctx, &result); err != nil {
		log.Println("repository, get all location, decoding", err)
		return nil, err
	}

	return result, nil
}

func (r *locationRepositoryImpl) GetLocationById(id *primitive.ObjectID) (*entity.Location, []*entity.Menu, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	var location entity.Location
	err := r.locationCollection.FindOne(ctx, filter).Decode(&location)
	if err != nil {
		log.Println("repository, get location by id", err)
		if err == mongo.ErrNoDocuments {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	menuFilter := bson.M{"location_id": id}
	ops := options.Find().SetProjection(bson.M{"steps": 0})
	cursor, err := r.menuCollection.Find(ctx, menuFilter, ops)
	if err != nil {
		log.Println("repository, get location by id, get menus", err)
		return nil, nil, err
	}

	var locationsMenus []*entity.Menu
	if err = cursor.All(ctx, &locationsMenus); err != nil {
		log.Println("repository, get location by id, get menus, decoding", err)
		return nil, nil, err
	}

	return &location, locationsMenus, nil
}

func (r *locationRepositoryImpl) GetMenuItmes(locationId, menuId *primitive.ObjectID) (*entity.Menu, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"location_id": locationId, "_id": menuId}
	var menu entity.Menu
	err := r.menuCollection.FindOne(ctx, filter).Decode(&menu)
	if err != nil {
		log.Println("repository, get menu items", err)
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &menu, nil
}
