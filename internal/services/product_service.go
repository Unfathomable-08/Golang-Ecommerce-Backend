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

type ProductsResult struct {
	Products     []models.Product `json:"products"`
	TotalPages   int              `json:"totalPages"`
	CurrentPage  int              `json:"currentPage"`
	TotalProducts int64           `json:"totalProducts"`
}

func GetProducts(page, limit int, category, search string, minPrice, maxPrice *float64) (*ProductsResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	if category != "" && category != "All" {
		filter["category"] = category
	}
	if search != "" {
		filter["$or"] = bson.A{
			bson.M{"name": bson.M{"$regex": search, "$options": "i"}},
			bson.M{"description": bson.M{"$regex": search, "$options": "i"}},
		}
	}
	if minPrice != nil || maxPrice != nil {
		priceFilter := bson.M{}
		if minPrice != nil { priceFilter["$gte"] = *minPrice }
		if maxPrice != nil { priceFilter["$lte"] = *maxPrice }
		filter["price"] = priceFilter
	}

	skip := int64((page - 1) * limit)
	opts := options.Find().
		SetSort(bson.M{"createdAt": -1}).
		SetSkip(skip).
		SetLimit(int64(limit))

	c := database.DB.Collection("products")
	total, err := c.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	cursor, err := c.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	if products == nil {
		products = []models.Product{}
	}

	return &ProductsResult{
		Products:     products,
		TotalPages:   int(math.Ceil(float64(total) / float64(limit))),
		CurrentPage:  page,
		TotalProducts: total,
	}, nil
}

func GetProduct(id string) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var product models.Product
	err = database.DB.Collection("products").FindOne(ctx, bson.M{"_id": oid}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func GetTrendingProducts(limit int) ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().
		SetSort(bson.D{{Key: "rating", Value: -1}, {Key: "reviews", Value: -1}}).
		SetLimit(int64(limit))

	cursor, err := database.DB.Collection("products").Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	if products == nil {
		products = []models.Product{}
	}
	return products, nil
}

func CreateProduct(p *models.Product) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	res, err := database.DB.Collection("products").InsertOne(ctx, p)
	if err != nil {
		return nil, err
	}
	p.ID = res.InsertedID.(bson.ObjectID)
	return p, nil
}

func UpdateProduct(id string, update bson.M) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	update["updatedAt"] = time.Now()

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var product models.Product
	err = database.DB.Collection("products").FindOneAndUpdate(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$set": update},
		opts,
	).Decode(&product)
	return &product, err
}

func DeleteProduct(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = database.DB.Collection("products").DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
