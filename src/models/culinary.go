package models

import (
	"context"
	"log"
	"time"

	"github.com/weeber-id/desatanjungbunga-backend/src/services"
	"github.com/weeber-id/desatanjungbunga-backend/src/tools"
	"github.com/weeber-id/desatanjungbunga-backend/src/variables"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Culinary collection model
type Culinary struct {
	BaseContent `bson:",inline"`

	Name  string `bson:"name" json:"name"`
	Image string `bson:"image" json:"image"`
	Slug  string `bsonn:"slug" json:"slug"`
	Price struct {
		Start string `bson:"start" json:"start"`
		End   string `bson:"end" json:"end"`
		Unit  string `bson:"unit" json:"unit"`
	} `bson:"price" json:"price"`
	OperationTime struct {
		From struct {
			Day  string `bson:"day" json:"day"`
			Time string `bson:"time" json:"time"`
		} `bson:"from" json:"from"`
		To struct {
			Day  string `bson:"day" json:"day"`
			Time string `bson:"time" json:"time"`
		} `bson:"to" json:"to"`
	} `bson:"operation_time" json:"operation_time"`
	Links []struct {
		Name string `bson:"name" json:"name"`
		Link string `bson:"link" json:"link"`
	} `bson:"links" json:"links"`
	ShortDescription string `bson:"short_description" json:"short_description"`
	Description      string `bson:"description" json:"description"`
}

// Collection pointer to this model
func (Culinary) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Culinary)
}

// Create new kuliner to database
func (k *Culinary) Create(ctx context.Context) error {
	k.CreatedAt = time.Now()
	k.UpdatedAt = time.Now()

	slug, err := tools.GenerateSlug(k.Name)
	if err != nil {
		log.Fatalf("error create slug culinary: %v", err)
	}
	k.Slug = slug

	result, err := k.Collection().InsertOne(ctx, *k)
	if err != nil {
		return err
	}

	k.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID read from database
func (k *Culinary) GetByID(ctx context.Context, id string) (found bool, err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	err = k.Collection().FindOne(ctx, bson.M{"_id": objectID}).Decode(k)
	return k.IsFoundFromError(err), err
}

// GetBySlug read from database
func (k *Culinary) GetBySlug(ctx context.Context, slug string) (found bool, err error) {
	err = k.Collection().FindOne(ctx, bson.M{"slug": slug}).Decode(k)
	return k.IsFoundFromError(err), err
}

// Update kuliner to database
func (k *Culinary) Update(ctx context.Context) error {
	k.UpdatedAt = time.Now()

	update := bson.M{"$set": *k}
	_, err := k.Collection().UpdateOne(ctx, bson.M{"_id": k.ID}, update)
	return err
}

// Delete kuliner from database
func (k *Culinary) Delete(ctx context.Context) error {
	return k.Collection().FindOneAndDelete(ctx, bson.M{"_id": k.ID}).Err()
}

// MultipleKuliner multiple model
type MultipleKuliner struct {
	baseList

	data []Culinary
}

// Collection kuliner mongo
func (MultipleKuliner) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Culinary)
}

// SortByName asc or desc
func (k *MultipleKuliner) SortByName(direction string) {
	numDirection := k.getDirectionFromStringToInt(direction)
	k.aggregate = append(k.aggregate, bson.M{
		"$sort": bson.M{"name": numDirection},
	})
}

// Get multiple kuliner from database
func (k *MultipleKuliner) Get(ctx context.Context) error {
	cur, err := k.Collection().Aggregate(ctx, k.aggregate)
	if err != nil {
		return err
	}

	for cur.Next(ctx) {
		var kuliner Culinary

		cur.Decode(&kuliner)
		k.data = append(k.data, kuliner)
	}

	return nil
}

// CountMaxPage execution
func (k *MultipleKuliner) CountMaxPage(ctx context.Context) uint {
	return k.countMaxPage(ctx, k.Collection())
}

// Data kuliner
func (k *MultipleKuliner) Data() []Culinary {
	return k.data
}
