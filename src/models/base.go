package models

import (
	"context"
	"log"
	"math"
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
	aggregate []bson.M

	Base      `bson:",inline"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

// baseList tools for multiple item
type baseList struct {
	aggregate []bson.M

	aggregateSearch []bson.M
	contentPerPage  int
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
	searchAggregate := []bson.M{
		{
			"$match": bson.M{"$text": bson.M{"$search": keyword}},
		},
		{
			"$sort": bson.M{"score": bson.M{"$meta": "textScore"}},
		},
	}

	l.aggregate = append(l.aggregate, searchAggregate...)
	l.aggregateSearch = append(l.aggregateSearch, searchAggregate...)
}

// FilterByPaginate aggregate
func (l *baseList) FilterByPaginate(page int, contentPerPage int) {
	l.contentPerPage = contentPerPage

	l.aggregate = append(l.aggregate, bson.M{
		"$skip": (page - 1) * contentPerPage,
	})

	l.aggregate = append(l.aggregate, bson.M{
		"$limit": contentPerPage,
	})
}

func (l *baseList) countDocuments(ctx context.Context, coll *mongo.Collection) uint {
	l.aggregate = append(l.aggregate, bson.M{
		"$group": bson.M{
			"_id":   nil,
			"count": bson.M{"$sum": 1},
		},
	})

	var result struct {
		Count uint `bson:"count"`
	}

	cur, err := coll.Aggregate(ctx, l.aggregate)
	if err != nil {
		return 0
	}

	for cur.Next(ctx) {
		cur.Decode(&result)
	}
	return result.Count
}

func (l *baseList) countMaxPage(ctx context.Context, coll *mongo.Collection) uint {
	if l.contentPerPage == 0 {
		l.contentPerPage = math.MaxInt32
	}

	l.aggregateSearch = append(l.aggregateSearch, bson.M{
		"$group": bson.M{
			"_id":   nil,
			"count": bson.M{"$sum": 1},
		},
	})
	l.aggregateSearch = append(l.aggregateSearch, bson.M{
		"$project": bson.M{
			"_id":      0,
			"max_page": bson.M{"$ceil": bson.M{"$divide": bson.A{"$count", l.contentPerPage}}},
		},
	})

	var result struct {
		MaxPage uint `bson:"max_page"`
	}

	cur, err := coll.Aggregate(ctx, l.aggregateSearch)
	if err != nil {
		return 0
	}

	for cur.Next(ctx) {
		cur.Decode(&result)
	}

	if result.MaxPage == 0 {
		return 1
	}
	return result.MaxPage
}
