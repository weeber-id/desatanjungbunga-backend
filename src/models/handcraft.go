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

// Handcraft collection model
type Handcraft struct {
	BaseContent `bson:",inline"`

	Name          string `bson:"name" json:"name"`
	Image         string `bson:"image" json:"image"`
	Price         string `bson:"price" json:"price"`
	Slug          string `bson:"slug" json:"slug"`
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
func (Handcraft) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Handcraft)
}

// Create new belanja to database
func (b *Handcraft) Create(ctx context.Context) error {
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()

	slug, err := tools.GenerateSlug(b.Name)
	if err != nil {
		log.Fatalf("error create slug article: %v", err)
	}
	b.Slug = slug

	result, err := b.Collection().InsertOne(ctx, *b)
	if err != nil {
		return err
	}

	b.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID read from database
func (b *Handcraft) GetByID(ctx context.Context, id string) (found bool, err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	err = b.Collection().FindOne(ctx, bson.M{"_id": objectID}).Decode(b)
	return b.IsFoundFromError(err), err
}

// GetBySlug read from database
func (b *Handcraft) GetBySlug(ctx context.Context, slug string) (found bool, err error) {
	err = b.Collection().FindOne(ctx, bson.M{"slug": slug}).Decode(b)
	return b.IsFoundFromError(err), err
}

// Update belanja to database
func (b *Handcraft) Update(ctx context.Context) error {
	b.UpdatedAt = time.Now()

	update := bson.M{"$set": *b}
	_, err := b.Collection().UpdateOne(ctx, bson.M{"_id": b.ID}, update)
	return err
}

// Delete belanja from database
func (b *Handcraft) Delete(ctx context.Context) error {
	return b.Collection().FindOneAndDelete(ctx, bson.M{"_id": b.ID}).Err()
}

// MultipleBelanja multiple model
type MultipleBelanja struct {
	baseList

	data []Handcraft
}

// Collection belanja mongo
func (MultipleBelanja) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Handcraft)
}

// SortByName asc or desc
func (b *MultipleBelanja) SortByName(direction string) {
	numDirection := b.getDirectionFromStringToInt(direction)
	b.aggregate = append(b.aggregate, bson.M{
		"$sort": bson.M{"name": numDirection},
	})
}

// Get multiple belanja from database
func (b *MultipleBelanja) Get(ctx context.Context) error {
	cur, err := b.Collection().Aggregate(ctx, b.aggregate)
	if err != nil {
		return err
	}

	for cur.Next(ctx) {
		var belanja Handcraft

		cur.Decode(&belanja)
		b.data = append(b.data, belanja)
	}

	return nil
}

// Data belanja
func (b *MultipleBelanja) Data() []Handcraft {
	return b.data
}
