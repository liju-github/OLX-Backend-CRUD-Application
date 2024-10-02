package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product represents a product entity in the application.
type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email       string             `bson:"email" json:"email"`       // Ensured valid email format
	Name        string             `bson:"name" json:"name" validate:"required"`               // Product name is required
	Description string             `bson:"description" json:"description" validate:"required"`               // Product name is required
	Category    string             `bson:"category" json:"category" validate:"required"`       // Category is required
	Price       float64            `bson:"price" json:"price" validate:"required,gte=0"`       // Price is required and must be greater than or equal to 0
	ImageURL    string             `bson:"image_url" json:"image_url" validate:"required,url"` // Image URL is required and must be a valid URL
	Address     string             `bson:"address" json:"address" validate:"required"`         // Address is required
	State       string             `bson:"state" json:"state" validate:"required"`             // State is required
	Pincode     string             `bson:"pincode" json:"pincode" validate:"required,len=6"`   // Pincode is required and must be exactly 6 characters long
}
