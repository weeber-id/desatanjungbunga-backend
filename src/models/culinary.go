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
		Monday struct {
			Open bool   `bson:"open" json:"open"`
			From string `bson:"from" json:"from"`
			To   string `bson:"to" json:"to"`
		} `bson:"monday" json:"monday"`
		Tuesday struct {
			Open bool   `bson:"open" json:"open"`
			From string `bson:"from" json:"from"`
			To   string `bson:"to" json:"to"`
		} `bson:"tuesday" json:"tuesday"`
		Wednesday struct {
			Open bool   `bson:"open" json:"open"`
			From string `bson:"from" json:"from"`
			To   string `bson:"to" json:"to"`
		} `bson:"wednesday" json:"wednesday"`
		Thursday struct {
			Open bool   `bson:"open" json:"open"`
			From string `bson:"from" json:"from"`
			To   string `bson:"to" json:"to"`
		} `bson:"thursday" json:"thursday"`
		Friday struct {
			Open bool   `bson:"open" json:"open"`
			From string `bson:"from" json:"from"`
			To   string `bson:"to" json:"to"`
		} `bson:"friday" json:"friday"`
		Saturday struct {
			Open bool   `bson:"open" json:"open"`
			From string `bson:"from" json:"from"`
			To   string `bson:"to" json:"to"`
		} `bson:"saturday" json:"saturday"`
		Sunday struct {
			Open bool   `bson:"open" json:"open"`
			From string `bson:"from" json:"from"`
			To   string `bson:"to" json:"to"`
		} `bson:"sunday" json:"sunday"`
	} `bson:"operation_time" json:"operation_time"`
	Links []struct {
		Name string `bson:"name" json:"name"`
		Link string `bson:"link" json:"link"`
	} `bson:"links" json:"links"`
	ShortDescription string `bson:"short_description" json:"short_description"`
	Description      string `bson:"description" json:"description"`
	Active           bool   `bson:"active" json:"active"`
	Recommendation   bool   `bson:"recommendation" json:"recommendation"`

	AuthorID primitive.ObjectID `bson:"author_id" json:"-"`
}

// Collection pointer to this model
func (Culinary) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Culinary)
}

// Create new kuliner to database
func (k *Culinary) Create(ctx context.Context, author *Admin) error {
	k.CreatedAt = time.Now()
	k.UpdatedAt = time.Now()
	k.AuthorID = author.ID
	k.Active = true

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

// FilterByAuthorID pipeline
func (k *MultipleKuliner) FilterByAuthorID(authorID string) *MultipleKuliner {
	objectID, _ := primitive.ObjectIDFromHex(authorID)

	filter := bson.M{
		"$match": bson.M{"author_id": objectID},
	}

	k.aggregateSearch = append(k.aggregateSearch, filter)
	k.aggregate = append(k.aggregate, filter)
	return k
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

// CountDocuments execution
func (k *MultipleKuliner) CountDocuments(ctx context.Context) uint {
	return k.countDocuments(ctx, k.Collection())
}

// CountMaxPage execution
func (k *MultipleKuliner) CountMaxPage(ctx context.Context) uint {
	return k.countMaxPage(ctx, k.Collection())
}

// Data kuliner
func (k *MultipleKuliner) Data() []Culinary {
	return k.data
}
