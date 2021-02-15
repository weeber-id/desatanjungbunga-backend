package models

import (
	"context"
	"time"

	"github.com/weeber-id/desatanjungbunga-backend/src/services"
	"github.com/weeber-id/desatanjungbunga-backend/src/tools"
	"github.com/weeber-id/desatanjungbunga-backend/src/variables"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Admin model database
type Admin struct {
	BaseContent `bson:",inline"`

	Name     string `bson:"name" json:"name"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"-"`
	Role     int    `bson:"role" json:"role"`
}

// Collection pointer to this model
func (Admin) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Admin)
}

// SetPassword for admin struct model
func (a *Admin) SetPassword(password string) {
	a.Password = tools.PasswordHashing(password, a.Username)
}

// IsPasswordMatch checker
func (a *Admin) IsPasswordMatch(password string) bool {
	hash := tools.PasswordHashing(password, a.Username)

	return a.Password == hash
}

// GetByUsername admin from database
func (a *Admin) GetByUsername(ctx context.Context, username string) (found bool, err error) {
	err = a.Collection().FindOne(ctx, bson.M{"username": username}).Decode(a)
	return a.IsFoundFromError(err), err
}

// Create new admin account
func (a *Admin) Create(ctx context.Context) error {
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()

	result, err := a.Collection().InsertOne(ctx, *a)
	if err != nil {
		return err
	}

	a.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// Update admin to database
func (a *Admin) Update(ctx context.Context) error {
	a.UpdatedAt = time.Now()

	update := bson.M{"$set": *a}
	_, err := a.Collection().UpdateOne(ctx, bson.M{"_id": a.ID}, update)
	return err
}

// Delete admin from database
func (a *Admin) Delete(ctx context.Context) error {
	return a.Collection().FindOneAndDelete(ctx, bson.M{"_id": a.ID}).Err()
}
