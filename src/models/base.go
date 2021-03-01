package models

import (
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Base model for collection in mongoDB
// suitable for logging collection
type Base struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}

// IsFoundFromError checker from scan mongo document
func (Base) IsFoundFromError(err error) bool {
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		log.Fatalf("Error in checking found from error: %v", err)
	}

	return true
}

// BaseContent model for collection in mongoDB
type BaseContent struct {
	Base      `bson:",inline"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

// baseList tools for multiple item
type baseList struct {
	aggregate []bson.M
}

func (baseList) getDirectionFromStringToInt(direction string) int {
	var numDirection int

	switch direction {
	case "asc":
		numDirection = 1
	default:
		numDirection = -1
	}

	return numDirection
}

// SortByDate asc (oldest) or desc (latest)
func (l *baseList) SortByDate(direction string) {
	numDirection := l.getDirectionFromStringToInt(direction)
	l.aggregate = append(l.aggregate, bson.M{
		"$sort": bson.M{"updated_at": numDirection},
	})
}

func (l *baseList) FilterBySearch(keyword string) {
	l.aggregate = append(l.aggregate, bson.M{
		"$match": bson.M{"$text": bson.M{"$search": keyword}},
	})

	l.aggregate = append(l.aggregate, bson.M{
		"$sort": bson.M{"score": bson.M{"$meta": "textScore"}},
	})
}

// FilterByPaginate aggregate
func (l *baseList) FilterByPaginate(page int, contentPerPage int) {
	l.aggregate = append(l.aggregate, bson.M{
		"$skip": (page - 1) * contentPerPage,
	})

	l.aggregate = append(l.aggregate, bson.M{
		"$limit": contentPerPage,
	})
}
