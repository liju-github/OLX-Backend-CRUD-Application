package repository

import (
	"context"
	"errors"

	"github.com/liju-github/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	Collection *mongo.Collection
}

func (repo *UserRepository) AddUser(user model.User) error {
	_, err := repo.Collection.InsertOne(context.TODO(), user)
	return err
}

func (repo *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := repo.Collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (repo *UserRepository) GetAllUsers() ([]model.User, error) {
    var users []model.User
    
    cursor, err := repo.Collection.Find(context.TODO(), bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.TODO()) 

    if err := cursor.All(context.TODO(), &users); err != nil {
        return nil, err
    }

    return users, nil
}


func (repo *UserRepository) UpdateUserImage(ctx context.Context, userEmail string, newImageUrl string) error {
    filter := bson.M{"email": userEmail}
    update := bson.M{"$set": bson.M{"image_url": newImageUrl}}

    _, err := repo.Collection.UpdateOne(ctx, filter, update)
    return err
}