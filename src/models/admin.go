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

const (
	AdminRoleSuperAdmin int = iota
	AdminRoleSeller
)

// Admin model database
type Admin struct {
	BaseContent `bson:",inline"`

	Name                string `bson:"name" json:"name"`
	Email               string `bson:"email" json:"email"`
	Address             string `bson:"address" json:"address"`
	DateofBirth         string `bson:"date_of_birth" json:"date_of_birth"`
	PhoneNumberWhatsapp string `bson:"phone_number_whatsapp" json:"phone_number_whatsapp"`
	Username            string `bson:"username" json:"username"`
	Password            string `bson:"password" json:"-"`
	Role                int    `bson:"role" json:"role"`
	ProfilePicture      string `bson:"profile_picture" json:"profile_picture"`
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

// GetByID admin from database
func (a *Admin) GetByID(ctx context.Context, id string) (found bool, err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	err = a.Collection().FindOne(ctx, bson.M{"_id": objectID}).Decode(a)
	return a.IsFoundFromError(err), err
}

func (a *Admin) GetByObjectID(ctx context.Context, id primitive.ObjectID) (found bool, err error) {
	err = a.Collection().FindOne(ctx, bson.M{"_id": id}).Decode(a)
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

// Admins model database
type Admins struct {
	baseList

	data []*Admin `bson:"data"`
}

// Collection pointer to this model
func (Admins) Collection() *mongo.Collection {
	return services.DB.Collection(variables.Collection.Admin)
}

// FilterByRole aggregation
func (a *Admins) FilterByRole(role int) *Admins {
	filter := bson.M{"$match": bson.M{"role": role}}

	a.aggregate = append(a.aggregate, filter)
	a.aggregateSearch = append(a.aggregateSearch, filter)
	return a
}

// Get list of admin from database
func (a *Admins) Get(ctx context.Context) error {
	cur, err := a.Collection().Aggregate(ctx, a.aggregate)
	if err != nil {
		return err
	}

	for cur.Next(ctx) {
		var admin Admin

		cur.Decode(&admin)
		a.data = append(a.data, &admin)
	}
	return nil
}

// CountMaxPage execution
func (a *Admins) CountMaxPage(ctx context.Context) uint {
	return a.countMaxPage(ctx, a.Collection())
}

// Data attributes
func (a *Admins) Data() []*Admin {
	return a.data
}
