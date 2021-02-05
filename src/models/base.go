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
	sort       bson.D
	pagination struct {
		skip  int64
		limit int64
	}
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

// SetPagination query
// page start from 1
func (l *baseList) SetPagination(page int64, contentPerPage int64) {
	l.pagination.skip = (page - 1) * contentPerPage
	l.pagination.limit = contentPerPage
}
