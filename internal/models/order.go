package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type OrderItem struct {
	ProductID bson.ObjectID `json:"productId" bson:"productId"`
	Name      string        `json:"name" bson:"name"`
	Price     float64       `json:"price" bson:"price"`
	Qty       int           `json:"qty" bson:"qty"`
}

type ShippingAddress struct {
	Address string `json:"address" bson:"address"`
	City    string `json:"city" bson:"city"`
	Zip     string `json:"zip" bson:"zip"`
	Country string `json:"country" bson:"country"`
}

type Order struct {
	ID              bson.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	User            string          `json:"user" bson:"user"` // ID or email of user, or guest ID
	Items           []OrderItem     `json:"items" bson:"items"`
	Total           float64         `json:"total" bson:"total"`
	Status          string          `json:"status" bson:"status"` // 'Pending', 'Processing', 'Shipped', 'Delivered', 'Cancelled'
	PaymentMethod   string          `json:"paymentMethod" bson:"paymentMethod"`
	ShippingAddress ShippingAddress `json:"shippingAddress" bson:"shippingAddress"`
	CreatedAt       time.Time       `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt" bson:"updatedAt"`
}
