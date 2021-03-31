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
func (Handcraft) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Handcraft)
}

// Create new belanja to database
func (b *Handcraft) Create(ctx context.Context, author *Admin) error {
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
	b.AuthorID = author.ID
	b.Active = true

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

func (b *Handcraft) GetByObjectID(ctx context.Context, objectID primitive.ObjectID) (found bool, err error) {
	err = b.Collection().FindOne(ctx, bson.M{"_id": objectID}).Decode(b)
	return b.IsFoundFromError(err), err
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

// FilterByAuthorID pipeline
func (b *MultipleBelanja) FilterByAuthorID(authorID string) *MultipleBelanja {
	objectID, _ := primitive.ObjectIDFromHex(authorID)

	filter := bson.M{
		"$match": bson.M{"author_id": objectID},
	}

	b.aggregateSearch = append(b.aggregateSearch, filter)
	b.aggregate = append(b.aggregate, filter)
	return b
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

// CountDocuments execution
func (b *MultipleBelanja) CountDocuments(ctx context.Context) uint {
	return b.countDocuments(ctx, b.Collection())
}

// CountMaxPage execution
func (b *MultipleBelanja) CountMaxPage(ctx context.Context) uint {
	return b.countMaxPage(ctx, b.Collection())
}

// Data belanja
func (b *MultipleBelanja) Data() []Handcraft {
	return b.data
}
