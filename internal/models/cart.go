package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CartItem struct {
	ProductID bson.ObjectID `json:"productId" bson:"productId"`
	Qty       int           `json:"qty" bson:"qty"`
}

type Cart struct {
	ID        bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CartID    string        `json:"cartId" bson:"cartId"`
	Items     []CartItem    `json:"items" bson:"items"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
}
