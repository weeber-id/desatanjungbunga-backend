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

// Wisata collection model
type Wisata struct {
	BaseContent `bson:",inline"`

	Title       string
	ImageCover  string `bson:"image_cover" json:"image_cover"`
	Author      string `bson:"author" json:"author"`
	Description string `bson:"description" json:"description"`
}

// Collection pointer to this model
func (Wisata) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Wisata)
}

// Create new wisata to database
func (w *Wisata) Create(ctx context.Context) error {
	w.CreatedAt = time.Now()
	w.UpdatedAt = time.Now()

	result, err := w.Collection().InsertOne(ctx, *w)
	if err != nil {
		return err
	}

	w.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID read from database
func (w *Wisata) GetByID(ctx context.Context, id string) (found bool, err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	err = w.Collection().FindOne(ctx, bson.M{"_id": objectID}).Decode(w)
	return w.IsFoundFromError(err), err
}

// Update wisata to database
func (w *Wisata) Update(ctx context.Context) error {
	w.UpdatedAt = time.Now()

	update := bson.M{"$set": *w}
	_, err := w.Collection().UpdateOne(ctx, bson.M{"_id": w.ID}, update)
	return err
}

// Delete wisata from database
func (w *Wisata) Delete(ctx context.Context) error {
	return w.Collection().FindOneAndDelete(ctx, bson.M{"_id": w.ID}).Err()
}

// MultipleWisata multiple model
type MultipleWisata struct {
	data []Wisata
}

// Collection wisata mongo
func (MultipleWisata) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Wisata)
}

// Get multiple wisata fomr database
func (w *MultipleWisata) Get(ctx context.Context) error {
	filter := bson.D{}

	opt := options.Find()
	opt.SetSort(bson.M{"_id": -1})

	cur, err := w.Collection().Find(ctx, filter, opt)
	if err != nil {
		return err
	}

	for cur.Next(ctx) {
		var wis Wisata

		cur.Decode(&wis)
		w.data = append(w.data, wis)
	}

	return nil
}

// Data wisata
func (w *MultipleWisata) Data() []Wisata {
	return w.data
}
