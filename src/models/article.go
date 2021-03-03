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

// Article collection model
type Article struct {
	BaseContent `bson:",inline"`

	Title      string `bson:"title" json:"title"`
	ImageCover string `bson:"image_cover" json:"image_cover"`
	Author     string `bson:"author" json:"author"`
	Body       string `bson:"body" json:"body"`
	Slug       string `bson:"slug" json:"slug"`
}

// Collection pointer to this model
func (Article) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Article)
}

// Create new article to database
func (a *Article) Create(ctx context.Context) error {
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	slug, err := tools.GenerateSlug(a.Title)
	if err != nil {
		log.Fatalf("error create slug article: %v", err)
	}
	a.Slug = slug

	result, err := a.Collection().InsertOne(ctx, *a)
	if err != nil {
		return err
	}

	a.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID read from database
func (a *Article) GetByID(ctx context.Context, id string) (found bool, err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	err = a.Collection().FindOne(ctx, bson.M{"_id": objectID}).Decode(a)
	return a.IsFoundFromError(err), err
}

// GetBySlug read from database
func (a *Article) GetBySlug(ctx context.Context, slug string) (found bool, err error) {
	err = a.Collection().FindOne(ctx, bson.M{"slug": slug}).Decode(a)
	return a.IsFoundFromError(err), err
}

// Update article to database
func (a *Article) Update(ctx context.Context) error {
	a.UpdatedAt = time.Now()

	update := bson.M{"$set": *a}
	_, err := a.Collection().UpdateOne(ctx, bson.M{"_id": a.ID}, update)
	return err
}

// Delete article from database
func (a *Article) Delete(ctx context.Context) error {
	return a.Collection().FindOneAndDelete(ctx, bson.M{"_id": a.ID}).Err()
}

// Articles multiple model
type Articles struct {
	baseList

	data []Article
}

// Collection article mongo
func (Articles) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Article)
}

// SortByTitle asc or desc
func (a *Articles) SortByTitle(direction string) {
	numDirection := a.getDirectionFromStringToInt(direction)
	a.aggregate = append(a.aggregate, bson.M{
		"$sort": bson.M{"title": numDirection},
	})
}

// Get multiple article from database
func (a *Articles) Get(ctx context.Context) error {
	cur, err := a.Collection().Aggregate(ctx, a.aggregate)
	if err != nil {
		return err
	}

	for cur.Next(ctx) {
		var art Article
		cur.Decode(&art)

		art.Body = ""
		a.data = append(a.data, art)
	}

	return nil
}

// CountMaxPage execution
func (a *Articles) CountMaxPage(ctx context.Context) uint {
	return a.countMaxPage(ctx, a.Collection())
}

// Data article
func (a *Articles) Data() []Article {
	return a.data
}
