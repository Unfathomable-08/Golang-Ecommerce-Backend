package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Review struct {
	ID        bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ProductID bson.ObjectID `json:"productId" bson:"productId"`
	User      string        `json:"user" bson:"user"` // Name or ID
	Rating    int           `json:"rating" bson:"rating"`
	Comment   string        `json:"comment" bson:"comment"`
	Status    string        `json:"status" bson:"status"` // 'Pending', 'Approved', 'Rejected'
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
}
