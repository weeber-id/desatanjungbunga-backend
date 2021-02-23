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

// Belanja collection model
type Belanja struct {
	BaseContent `bson:",inline"`

	Title       string
	ImageCover  string `bson:"image_cover" json:"image_cover"`
	Description string `bson:"description" json:"description"`
}

// Collection pointer to this model
func (Belanja) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Handcraft)
}

// Create new belanja to database
func (b *Belanja) Create(ctx context.Context) error {
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()

	result, err := b.Collection().InsertOne(ctx, *b)
	if err != nil {
		return err
	}

	b.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID read from database
func (b *Belanja) GetByID(ctx context.Context, id string) (found bool, err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	err = b.Collection().FindOne(ctx, bson.M{"_id": objectID}).Decode(b)
	return b.IsFoundFromError(err), err
}

// Update belanja to database
func (b *Belanja) Update(ctx context.Context) error {
	b.UpdatedAt = time.Now()

	update := bson.M{"$set": *b}
	_, err := b.Collection().UpdateOne(ctx, bson.M{"_id": b.ID}, update)
	return err
}

// Delete belanja from database
func (b *Belanja) Delete(ctx context.Context) error {
	return b.Collection().FindOneAndDelete(ctx, bson.M{"_id": b.ID}).Err()
}

// MultipleBelanja multiple model
type MultipleBelanja struct {
	data []Belanja
}

// Collection belanja mongo
func (MultipleBelanja) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Handcraft)
}

// Get multiple belanja from database
func (b *MultipleBelanja) Get(ctx context.Context) error {
	filter := bson.D{}

	opt := options.Find()
	opt.SetSort(bson.M{"_id": -1})

	cur, err := b.Collection().Find(ctx, filter, opt)
	if err != nil {
		return err
	}

	for cur.Next(ctx) {
		var belanja Belanja

		cur.Decode(&belanja)
		b.data = append(b.data, belanja)
	}

	return nil
}

// Data belanja
func (b *MultipleBelanja) Data() []Belanja {
	return b.data
}
