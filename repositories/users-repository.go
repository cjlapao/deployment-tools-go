package repositories

import (
	"context"
	"time"

	"github.com/cjlapao/deployment-tools-go/entities"

	"go.mongodb.org/mongo-driver/bson"
)

// GetUserByEmail Gets a record by ID
func (a *Repository) GetUserByEmail(email string) entities.User {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user entities.User

	filter := bson.D{{Key: "email", Value: email}}
	collection := a.Factory.Database.Collection(UsersCollectionName)

	collection.FindOne(ctx, filter).Decode(&user)

	return user
}
