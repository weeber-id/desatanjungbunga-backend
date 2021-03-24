package models

import (
	"context"
	"time"

	"github.com/weeber-id/desatanjungbunga-backend/src/services"
	"github.com/weeber-id/desatanjungbunga-backend/src/variables"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type About struct {
	BaseContent `bson:",inline"`

	Name           string `bson:"name" json:"name"`
	ProfilePicture string `bson:"profile_picture" json:"profile_picture"`
	Position       string `bson:"position" json:"position"`
	Body           string `bson:"body" json:"body"`
}

func (About) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.About)
}

func (a *About) Get(ctx context.Context) (found bool, err error) {
	err = a.Collection().FindOne(ctx, bson.M{}).Decode(a)
	return a.IsFoundFromError(err), err
}

func (a *About) Create(ctx context.Context) error {
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	result, err := a.Collection().InsertOne(ctx, *a)
	if err != nil {
		return err
	}

	a.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (a *About) Update(ctx context.Context) error {
	a.UpdatedAt = time.Now()

	update := bson.M{"$set": *a}
	_, err := a.Collection().UpdateOne(ctx, bson.M{"_id": a.ID}, update)
	return err
}

func (a *About) Delete(ctx context.Context) error {
	return a.Collection().FindOneAndDelete(ctx, bson.M{"_id": a.ID}).Err()
}
