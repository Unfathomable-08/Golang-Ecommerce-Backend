package services

import (
	"context"
	"time"

	"api/internal/database"
	"api/internal/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func GetCarts() ([]models.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.DB.Collection("carts").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var carts []models.Cart
	if err := cursor.All(ctx, &carts); err != nil {
		return nil, err
	}
	if carts == nil {
		carts = []models.Cart{}
	}
	return carts, nil
}

func UpsertCart(cartID string, productID bson.ObjectID, qty int) (*models.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := database.DB.Collection("carts")

	// Try to find existing cart
	var cart models.Cart
	err := c.FindOne(ctx, bson.M{"cartId": cartID}).Decode(&cart)

	if err != nil {
		// Cart doesn't exist — create it
		newCart := models.Cart{
			CartID:    cartID,
			Items:     []models.CartItem{{ProductID: productID, Qty: qty}},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		res, err := c.InsertOne(ctx, newCart)
		if err != nil {
			return nil, err
		}
		newCart.ID = res.InsertedID.(bson.ObjectID)
		return &newCart, nil
	}

	// Cart exists — check if item already in cart
	itemExists := false
	for i, item := range cart.Items {
		if item.ProductID == productID {
			cart.Items[i].Qty += qty
			itemExists = true
			break
		}
	}
	if !itemExists {
		cart.Items = append(cart.Items, models.CartItem{ProductID: productID, Qty: qty})
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Cart
	err = c.FindOneAndUpdate(
		ctx,
		bson.M{"cartId": cartID},
		bson.M{"$set": bson.M{"items": cart.Items, "updatedAt": time.Now()}},
		opts,
	).Decode(&updated)
	return &updated, err
}
