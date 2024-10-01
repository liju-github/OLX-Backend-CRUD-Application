package repository

import (
	"context"

	"github.com/liju-github/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository struct {
	Collection *mongo.Collection
}

func (repo *ProductRepository) AddProduct(product model.Product) error {
	_, err := repo.Collection.InsertOne(context.TODO(), product)
	return err
}

func (repo *ProductRepository) GetProductByID(id string) (*model.Product, error) {
	var product model.Product
	objectID, _ := primitive.ObjectIDFromHex(id)
	err := repo.Collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&product)
	return &product, err
}

func (repo *ProductRepository) GetAllProducts() ([]model.Product, error) {
	var products []model.Product

	cursor, err := repo.Collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (repo *ProductRepository) GetAllProductsByUserEmail(email string) ([]model.Product, error) {
	var products []model.Product

	// Find products where UserEmail matches the provided email
	cursor, err := repo.Collection.Find(context.TODO(), bson.M{"email": email})
	if err != nil {
		return nil, err // Return the error if the query fails
	}
	defer cursor.Close(context.TODO())

	// Decode each product found into the products slice
	for cursor.Next(context.TODO()) {
		var product model.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err // Return the error if decoding fails
		}
		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		return nil, err // Return any error encountered during iteration
	}

	return products, nil // Return the list of products found
}
