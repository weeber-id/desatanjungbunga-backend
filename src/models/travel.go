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

// Travel collection model
type Travel struct {
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
	ShortDescription string `bson:"short_description" json:"short_description"`
	Description      string `bson:"description" json:"description"`
	Active           bool   `bson:"active" json:"active"`
	Recommendation   bool   `bson:"recommendation" json:"recommendation"`

	AuthorID primitive.ObjectID `bson:"author_id" json:"-"`
}

// Collection pointer to this model
func (Travel) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Travel)
}

// Create new wisata to database
func (w *Travel) Create(ctx context.Context, author *Admin) error {
	w.CreatedAt = time.Now()
	w.UpdatedAt = time.Now()
	w.AuthorID = author.ID
	w.Active = true

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

func (w *Travel) GetByObjectID(ctx context.Context, objectID primitive.ObjectID) (found bool, err error) {
	err = w.Collection().FindOne(ctx, bson.M{"_id": objectID}).Decode(w)
	return w.IsFoundFromError(err), err
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
	baseList

	data []Travel
}

// Collection wisata mongo
func (MultipleWisata) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Travel)
}

// SortByName asc or desc
func (w *MultipleWisata) SortByName(direction string) {
	numDirection := w.getDirectionFromStringToInt(direction)
	w.aggregate = append(w.aggregate, bson.M{
		"$sort": bson.M{"name": numDirection},
	})
}

// FilterByAuthorID pipeline
func (w *MultipleWisata) FilterByAuthorID(authorID string) *MultipleWisata {
	objectID, _ := primitive.ObjectIDFromHex(authorID)

	filter := bson.M{
		"$match": bson.M{"author_id": objectID},
	}

	w.aggregateSearch = append(w.aggregateSearch, filter)
	w.aggregate = append(w.aggregate, filter)
	return w
}

// Get multiple wisata from database
func (w *MultipleWisata) Get(ctx context.Context) error {
	cur, err := w.Collection().Aggregate(ctx, w.aggregate)
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

// CountDocuments execution
func (w *MultipleWisata) CountDocuments(ctx context.Context) uint {
	return w.countDocuments(ctx, w.Collection())
}

// CountMaxPage execution
func (w *MultipleWisata) CountMaxPage(ctx context.Context) uint {
	return w.countMaxPage(ctx, w.Collection())
}

// Data wisata
func (w *MultipleWisata) Data() []Travel {
	return w.data
}
