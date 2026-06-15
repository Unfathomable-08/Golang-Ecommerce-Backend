package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Image struct {
	URL string `json:"url" bson:"url"`
	ID  string `json:"id" bson:"id"`
}

type Color struct {
	Name string `json:"name" bson:"name"`
	Hex  string `json:"hex" bson:"hex"`
}

type Product struct {
	ID            bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name          string        `json:"name" bson:"name"`
	Description   string        `json:"description" bson:"description"`
	Price         float64       `json:"price" bson:"price"`
	OriginalPrice *float64      `json:"originalPrice,omitempty" bson:"originalPrice,omitempty"`
	Category      string        `json:"category" bson:"category"`
	InStock       bool          `json:"inStock" bson:"inStock"`
	Featured      bool          `json:"featured" bson:"featured"`
	Images        []Image       `json:"images" bson:"images"`
	Colors        []Color       `json:"colors" bson:"colors"`
	Sizes         []string      `json:"sizes" bson:"sizes"`
	Rating        float64       `json:"rating" bson:"rating"`
	Reviews       int           `json:"reviews" bson:"reviews"`
	CreatedAt     time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt" bson:"updatedAt"`
}
