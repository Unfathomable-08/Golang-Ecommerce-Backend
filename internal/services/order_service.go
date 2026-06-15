package services

import (
	"context"
	"math"
	"time"

	"api/internal/database"
	"api/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type OrdersResult struct {
	Orders      []models.Order `json:"orders"`
	TotalPages  int            `json:"totalPages"`
	CurrentPage int            `json:"currentPage"`
	TotalOrders int64          `json:"totalOrders"`
}

func GetOrders(page, limit int) (*OrdersResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	skip := int64((page - 1) * limit)
	opts := options.Find().
		SetSort(bson.M{"createdAt": -1}).
		SetSkip(skip).
		SetLimit(int64(limit))

	c := database.DB.Collection("orders")
	total, err := c.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	cursor, err := c.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orders []models.Order
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	if orders == nil {
		orders = []models.Order{}
	}

	return &OrdersResult{
		Orders:      orders,
		TotalPages:  int(math.Ceil(float64(total) / float64(limit))),
		CurrentPage: page,
		TotalOrders: total,
	}, nil
}

func GetOrder(id string) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var order models.Order
	err = database.DB.Collection("orders").FindOne(ctx, bson.M{"_id": oid}).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func CreateOrder(o *models.Order) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	o.Status = "Pending"
	if o.PaymentMethod == "" {
		o.PaymentMethod = "COD"
	}
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()

	res, err := database.DB.Collection("orders").InsertOne(ctx, o)
	if err != nil {
		return nil, err
	}
	o.ID = res.InsertedID.(bson.ObjectID)
	return o, nil
}

func UpdateOrder(id string, update bson.M) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	update["updatedAt"] = time.Now()

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var order models.Order
	err = database.DB.Collection("orders").FindOneAndUpdate(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$set": update},
		opts,
	).Decode(&order)
	return &order, err
}

func CancelOrder(id string) (*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var current models.Order
	err = database.DB.Collection("orders").FindOne(ctx, bson.M{"_id": oid}).Decode(&current)
	if err != nil {
		return nil, err
	}
	if current.Status == "Shipped" || current.Status == "Delivered" {
		return nil, ErrCannotCancel
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var order models.Order
	err = database.DB.Collection("orders").FindOneAndUpdate(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{"status": "Cancelled", "updatedAt": time.Now()}},
		opts,
	).Decode(&order)
	return &order, err
}
