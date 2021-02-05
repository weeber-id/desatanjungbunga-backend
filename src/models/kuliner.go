package models

import (
	"context"
	"time"

	"github.com/weeber-id/desatanjungbunga-backend/src/services"
	"github.com/weeber-id/desatanjungbunga-backend/src/variables"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Kuliner collection model
type Kuliner struct {
	BaseContent `bson:",inline"`

	Title       string
	ImageCover  string `bson:"image_cover" json:"image_cover"`
	Description string `bson:"description" json:"description"`
}

// Collection pointer to this model
func (Kuliner) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Kuliner)
}

// Create new kuliner to database
func (k *Kuliner) Create(ctx context.Context) error {
	k.CreatedAt = time.Now()
	k.UpdatedAt = time.Now()

	result, err := k.Collection().InsertOne(ctx, *k)
	if err != nil {
		return err
	}

	k.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID read from database
func (k *Kuliner) GetByID(ctx context.Context, id string) (found bool, err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	err = k.Collection().FindOne(ctx, bson.M{"_id": objectID}).Decode(k)
	return k.IsFoundFromError(err), err
}

// Update kuliner to database
func (k *Kuliner) Update(ctx context.Context) error {
	k.UpdatedAt = time.Now()

	update := bson.M{"$set": *k}
	_, err := k.Collection().UpdateOne(ctx, bson.M{"_id": k.ID}, update)
	return err
}

// Delete kuliner from database
func (k *Kuliner) Delete(ctx context.Context) error {
	return k.Collection().FindOneAndDelete(ctx, bson.M{"_id": k.ID}).Err()
}

// MultipleKuliner multiple model
type MultipleKuliner struct {
	data []Kuliner
}

// Collection kuliner mongo
func (MultipleKuliner) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Kuliner)
}

// Get multiple kuliner fomr database
func (k *MultipleKuliner) Get(ctx context.Context) error {
	filter := bson.D{}

	opt := options.Find()
	opt.SetSort(bson.M{"_id": -1})

	cur, err := k.Collection().Find(ctx, filter, opt)
	if err != nil {
		return err
	}

	for cur.Next(ctx) {
		var kuliner Kuliner

		cur.Decode(&kuliner)
		k.data = append(k.data, kuliner)
	}

	return nil
}

// Data kuliner
func (k *MultipleKuliner) Data() []Kuliner {
	return k.data
}
