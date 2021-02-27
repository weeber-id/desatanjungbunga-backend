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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Travel collection model
type Travel struct {
	BaseContent   `bson:",inline"`
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
	ShortDescription string `bson:"short_description" json:"short_description"`
	Description      string `bson:"description" json:"description"`
}

// Collection pointer to this model
func (Travel) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Travel)
}

// Create new wisata to database
func (w *Travel) Create(ctx context.Context) error {
	w.CreatedAt = time.Now()
	w.UpdatedAt = time.Now()

	slug, err := tools.GenerateSlug(w.Name)
	if err != nil {
		log.Fatalf("error create slug article: %v", err)
	}
	w.Slug = slug

	result, err := w.Collection().InsertOne(ctx, *w)
	if err != nil {
		return err
	}

	w.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID read from database
func (w *Travel) GetByID(ctx context.Context, id string) (found bool, err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	err = w.Collection().FindOne(ctx, bson.M{"_id": objectID}).Decode(w)
	return w.IsFoundFromError(err), err
}

// GetBySlug read from database
func (w *Travel) GetBySlug(ctx context.Context, slug string) (found bool, err error) {
	err = w.Collection().FindOne(ctx, bson.M{"slug": slug}).Decode(w)
	return w.IsFoundFromError(err), err
}

// Update wisata to database
func (w *Travel) Update(ctx context.Context) error {
	w.UpdatedAt = time.Now()

	update := bson.M{"$set": *w}
	_, err := w.Collection().UpdateOne(ctx, bson.M{"_id": w.ID}, update)
	return err
}

// Delete wisata from database
func (w *Travel) Delete(ctx context.Context) error {
	return w.Collection().FindOneAndDelete(ctx, bson.M{"_id": w.ID}).Err()
}

// MultipleWisata multiple model
type MultipleWisata struct {
	data []Travel
}

// Collection wisata mongo
func (MultipleWisata) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Travel)
}

// Get multiple wisata from database
func (w *MultipleWisata) Get(ctx context.Context) error {
	filter := bson.D{}

	opt := options.Find()
	opt.SetSort(bson.M{"_id": -1})

	cur, err := w.Collection().Find(ctx, filter, opt)
	if err != nil {
		return err
	}

	for cur.Next(ctx) {
		var wis Travel

		cur.Decode(&wis)
		w.data = append(w.data, wis)
	}

	return nil
}

// Data wisata
func (w *MultipleWisata) Data() []Travel {
	return w.data
}
