package usermanagement

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/litmuschaos/litmus/litmus-portal/graphql-server/pkg/database/mongodb"
)

// CreateUser inserts a new user to the database
func CreateUser(ctx context.Context, user *User) error {
	err := mongodb.Operator.Create(ctx, mongodb.UserCollection, user)
	if err != nil {
		log.Print("Error creating User : ", err)
		return err
	}

	return nil
}

// GetUserByUserName returns user details based on username
func GetUserByUserName(ctx context.Context, username string) (*User, error) {
	var user = new(User)
	query := bson.D{{"username", username}}

	result, err := mongodb.Operator.Get(ctx, mongodb.UserCollection, query)
	if err != nil {
		log.Print("Error getting user with username: ", username, "\nError message: ", err)
		return nil, err
	}
	err = result.Decode(user)
	if err != nil {
		log.Print("Error unmarshalling the result in user struct: ", err)
		return nil, err
	}

	return user, err
}

// GetUserByUserID returns user details based on userID
func GetUserByUserID(ctx context.Context, userID string) (*User, error) {
	var user = new(User)
	query := bson.D{{"_id", userID}}
	result, err := mongodb.Operator.Get(ctx, mongodb.UserCollection, query)
	if err != nil {
		log.Print("Error getting user with userID: ", userID, "\nError message: ", err)
		return nil, err
	}
	err = result.Decode(user)
	if err != nil {
		log.Print("Error unmarshalling the result in user struct ", err)
		return nil, err
	}

	return user, err
}

// GetUsers returns the list of users present in the project
func GetUsers(ctx context.Context) ([]User, error) {
	query := bson.D{{}}
	result, err := mongodb.Operator.List(ctx, mongodb.UserCollection, query)
	if err != nil {
		log.Print("Error getting users : ", err)
		return []User{}, err
	}
	var users []User
	err = result.All(ctx, &users)
	if err != nil {
		log.Print("Error unmarshalling the result in users array: ", err)
		return []User{}, err
	}

	return users, nil
}

// UpdateUser updates the details of user in both user and project DB collections
func UpdateUser(ctx context.Context, user *User) error {

	filter := bson.D{{"_id", user.ID}}
	update := bson.D{
		{"$set", bson.D{
			{"name", user.Name},
			{"email", user.Email},
			{"company_name", user.CompanyName},
			{"updated_at", user.UpdatedAt},
		}},
	}

	_, err := mongodb.Operator.Update(ctx, mongodb.UserCollection, filter, update)
	if err != nil {
		log.Print("Error updating user: ", err)
		return err
	}

	opts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{"elem.user_id", user.ID}},
		},
	})
	filter = bson.D{{}}
	update = bson.D{
		{"$set", bson.D{
			{"members.$[elem].name", user.Name},
			{"members.$[elem].email", user.Email},
			{"members.$[elem].company_name", user.CompanyName},
		}},
	}

	_, err = mongodb.Operator.UpdateMany(ctx, mongodb.ProjectCollection, filter, update, opts)
	if err != nil {
		log.Print("Error updating user in projects : ", err)
		return err
	}

	return nil
}
