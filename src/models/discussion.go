package models

import (
	"context"
	"errors"
	"time"

	"github.com/weeber-id/desatanjungbunga-backend/src/services"
	"github.com/weeber-id/desatanjungbunga-backend/src/tools"
	"github.com/weeber-id/desatanjungbunga-backend/src/variables"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Discussion collection model
type Discussion struct {
	BaseContent `bson:",inline"`

	ParentID    primitive.ObjectID `bson:"parent_id,omitempty" json:"-"`
	Name        string             `bson:"name" json:"name"`
	Email       string             `bson:"email" json:"email"`
	Body        string             `bson:"body" json:"body"`
	ContentName string             `bson:"content_name" json:"content_name"`
	ContentID   primitive.ObjectID `bson:"content_id" json:"content_id"`

	Questions []Discussion `bson:"-" json:"questions"`
}

// Collection pointor to this model
func (Discussion) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Discussion)
}

// SetParentID of this discussion answer the question
func (d *Discussion) SetParentID(parentID string) (found bool, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	question := new(Discussion)
	found, err = question.GetByID(ctx, parentID)
	if !found {
		return found, err
	}

	d.ParentID, err = primitive.ObjectIDFromHex(parentID)
	if err != nil {
		return true, err
	}
	return true, nil
}

func (d *Discussion) SetContentNameAndID(contentName string, contentID string) error {
	allowedContentName := []string{
		variables.Collection.Article,
		variables.Collection.Travel,
		variables.Collection.Culinary,
		variables.Collection.Handcraft,
		variables.Collection.Lodging,
	}
	if !tools.InArrayStrings(allowedContentName, contentName) {
		return errors.New("invalid content name")
	}
	d.ContentName = contentName

	objectContentID, err := primitive.ObjectIDFromHex(contentID)
	if err != nil {
		return err
	}
	d.ContentID = objectContentID
	return nil
}

// Create new discussion to database
func (d *Discussion) Create(ctx context.Context) error {
	d.CreatedAt = time.Now()
	d.UpdatedAt = time.Now()

	result, err := d.Collection().InsertOne(ctx, *d)
	if err != nil {
		return err
	}

	d.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID read from database
func (d *Discussion) GetByID(ctx context.Context, id string) (found bool, err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	err = d.Collection().FindOne(ctx, bson.M{"_id": objectID}).Decode(d)
	return d.IsFoundFromError(err), err
}

// GetQuestions for this discussion
func (d *Discussion) GetQuestions(ctx context.Context) error {
	cur, err := d.Collection().Find(ctx, bson.M{"parent_id": d.ID}, options.Find())
	if err != nil {
		return err
	}

	for cur.Next(ctx) {
		var question Discussion
		cur.Decode(&question)
		d.Questions = append(d.Questions, question)
	}
	return nil
}

// Update discussion to database
func (d *Discussion) Update(ctx context.Context) error {
	d.UpdatedAt = time.Now()

	update := bson.M{"$set": *d}
	_, err := d.Collection().UpdateOne(ctx, bson.M{"_id": d.ID}, update)
	return err
}

// Delete discussion from database
func (d *Discussion) Delete(ctx context.Context) error {
	return d.Collection().FindOneAndDelete(ctx, bson.M{"_id": d.ID}).Err()
}

// Discussions multiple model
type Discussions struct {
	baseList

	data []Discussion
}

// Collection discussion mongo
func (Discussions) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Discussion)
}

// FilterByContentNameID query
func (d *Discussions) FilterByContentNameID(contentName string, contentID string) error {
	objectID, err := primitive.ObjectIDFromHex(contentID)
	if err != nil {
		return err
	}

	allowedContentName := []string{
		variables.Collection.Article,
		variables.Collection.Travel,
		variables.Collection.Culinary,
		variables.Collection.Handcraft,
		variables.Collection.Lodging,
	}

	if !tools.InArrayStrings(allowedContentName, contentName) {
		return errors.New("invalid content name")
	}

	d.aggregate = append(d.aggregate, bson.M{
		"$match": bson.D{
			{Key: "content_name", Value: contentName},
			{Key: "content_id", Value: objectID},
		},
	})

	return nil
}

// FilterOnlyQuestion query
func (d *Discussions) FilterOnlyQuestion() {
	d.aggregate = append(d.aggregate, bson.M{
		"$match": bson.M{"parent_id": nil},
	})
}

// FilterOnlyAnswer query
func (d *Discussions) FilterOnlyAnswer(parentID string) {
	parentObjectID, _ := primitive.ObjectIDFromHex(parentID)

	d.aggregate = append(d.aggregate, bson.M{
		"$match": bson.M{"parent_id": parentObjectID},
	})
}

// Get from discussion from database
func (d *Discussions) Get(ctx context.Context, showQuestion bool) error {
	cur, err := d.Collection().Aggregate(ctx, d.aggregate)
	if err != nil {
		return err
	}

	for cur.Next(ctx) {
		var question Discussion
		cur.Decode(&question)

		if showQuestion {
			err := question.GetQuestions(ctx)
			if err != nil {
				return err
			}
		}

		d.data = append(d.data, question)
	}

	return nil
}

// Data discussions
func (d *Discussions) Data() []Discussion {
	return d.data
}
