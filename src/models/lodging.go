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

// Lodging structure model
type Lodging struct {
	BaseContent `bson:",inline"`

	Name  string `bson:"name" json:"name"`
	Image string `bson:"image" json:"image"`
	Slug  string `bson:"slug" json:"slug"`
	Price struct {
		Value string `bson:"value" json:"value"`
		Unit  string `bson:"unit" json:"unit"`
	} `bson:"price" json:"price"`
	OperationTime string `bson:"operation_time" json:"operation_time"`
	Links         []struct {
		Name string `bson:"name" json:"name"`
		Link string `bson:"link" json:"link"`
	} `bson:"links" json:"links"`
	FacilitiesID     []string `bson:"facilities_id" json:"-"`
	ShortDescription string   `bson:"short_description" json:"short_description"`
	Description      string   `bson:"description" json:"description"`

	// Custom fields
	Facilities []LodgingFacility `bson:"-" json:"facilities"`
}

// Collection pointer to this model
func (Lodging) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Lodging)
}

// Create new lodging to database
func (l *Lodging) Create(ctx context.Context) error {
	l.CreatedAt = time.Now()
	l.UpdatedAt = time.Now()

	slug, err := tools.GenerateSlug(l.Name)
	if err != nil {
		log.Fatalf("error create slug article: %v", err)
	}
	l.Slug = slug

	result, err := l.Collection().InsertOne(ctx, *l)
	if err != nil {
		return err
	}

	l.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID read from database
func (l *Lodging) GetByID(ctx context.Context, id string) (found bool, err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	err = l.Collection().FindOne(ctx, bson.M{"_id": objectID}).Decode(l)
	return l.IsFoundFromError(err), err
}

// GetBySlug read from database
func (l *Lodging) GetBySlug(ctx context.Context, slug string) (found bool, err error) {
	err = l.Collection().FindOne(ctx, bson.M{"slug": slug}).Decode(l)
	return l.IsFoundFromError(err), err
}

// LoadFacilitiesDetail after query
func (l *Lodging) LoadFacilitiesDetail(ctx context.Context) error {
	facilities := new(MultipleLodgingFacility)
	if err := facilities.GetByIDs(ctx, l.FacilitiesID); err != nil {
		return err
	}

	l.Facilities = facilities.Data()
	return nil
}

// Update lodging to database
func (l *Lodging) Update(ctx context.Context) error {
	l.UpdatedAt = time.Now()

	update := bson.M{"$set": *l}
	_, err := l.Collection().UpdateOne(ctx, bson.M{"_id": l.ID}, update)
	return err
}

// Delete lodging from database
func (l *Lodging) Delete(ctx context.Context) error {
	return l.Collection().FindOneAndDelete(ctx, bson.M{"_id": l.ID}).Err()
}

// MultipleLodging collection structure
type MultipleLodging struct {
	baseList

	data []Lodging
}

// Collection pointer to this model
func (MultipleLodging) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Lodging)
}

// SortByName asc or desc
func (m *MultipleLodging) SortByName(direction string) {
	numDirection := m.getDirectionFromStringToInt(direction)
	m.aggregate = append(m.aggregate, bson.M{
		"$sort": bson.M{"name": numDirection},
	})
}

// Get multiple lodging from database
func (m *MultipleLodging) Get(ctx context.Context) error {
	cur, err := m.Collection().Aggregate(ctx, m.aggregate)
	if err != nil {
		return err
	}

	for cur.Next(ctx) {
		var lodging Lodging

		cur.Decode(&lodging)
		m.data = append(m.data, lodging)
	}

	return nil
}

// CountMaxPage execution
func (m *MultipleLodging) CountMaxPage(ctx context.Context) uint {
	return m.countMaxPage(ctx, m.Collection())
}

// Data lodging
func (m *MultipleLodging) Data() []Lodging {
	return m.data
}

// LodgingFacility collection model
type LodgingFacility struct {
	BaseContent `bson:",inline"`

	Name string `bson:"name" json:"name"`
	Icon string `bson:"icon" json:"icon"`
}

// Collection pointer to this model
func (LodgingFacility) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.LodgindFacilities)
}

// Create new facility to database
func (l *LodgingFacility) Create(ctx context.Context) error {
	l.CreatedAt = time.Now()
	l.UpdatedAt = time.Now()

	result, err := l.Collection().InsertOne(ctx, *l)
	if err != nil {
		return err
	}

	l.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID read from database
func (l *LodgingFacility) GetByID(ctx context.Context, id string) (found bool, err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	err = l.Collection().FindOne(ctx, bson.M{"_id": objectID}).Decode(l)
	return l.IsFoundFromError(err), err
}

// Update facility to database
func (l *LodgingFacility) Update(ctx context.Context) error {
	l.UpdatedAt = time.Now()

	update := bson.M{"$set": *l}
	_, err := l.Collection().UpdateOne(ctx, bson.M{"_id": l.ID}, update)
	return err
}

// Delete facility from database
func (l *LodgingFacility) Delete(ctx context.Context) error {
	return l.Collection().FindOneAndDelete(ctx, bson.M{"_id": l.ID}).Err()
}

// MultipleLodgingFacility model
type MultipleLodgingFacility struct {
	data []LodgingFacility
}

// Collection this struct
func (MultipleLodgingFacility) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.LodgindFacilities)
}

// GetByIDs multiple facility from database
func (l *MultipleLodgingFacility) GetByIDs(ctx context.Context, ids []string) error {
	var objectIDs []primitive.ObjectID

	for _, idstring := range ids {
		objectID, err := primitive.ObjectIDFromHex(idstring)
		if err != nil {
			continue
		}
		objectIDs = append(objectIDs, objectID)
	}

	cur, err := l.Collection().Find(ctx, bson.M{
		"_id": bson.M{"$in": objectIDs},
	})
	if err != nil {
		return err
	}

	for cur.Next(ctx) {
		var facility LodgingFacility

		cur.Decode(&facility)
		l.data = append(l.data, facility)
	}

	return nil
}

// Get all execution query
func (l *MultipleLodgingFacility) Get(ctx context.Context) error {
	cur, err := l.Collection().Find(ctx, bson.M{})
	if err != nil {
		return err
	}

	for cur.Next(ctx) {
		var facility LodgingFacility

		cur.Decode(&facility)
		l.data = append(l.data, facility)
	}

	return nil
}

// Data property
func (l *MultipleLodgingFacility) Data() []LodgingFacility {
	return l.data
}
