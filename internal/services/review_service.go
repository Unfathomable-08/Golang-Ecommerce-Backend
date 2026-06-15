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

type ReviewsResult struct {
	Reviews      []models.Review `json:"reviews"`
	TotalPages   int             `json:"totalPages"`
	CurrentPage  int             `json:"currentPage"`
	TotalReviews int64           `json:"totalReviews"`
}

func GetReviews(page, limit int, status string) (*ReviewsResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}

	skip := int64((page - 1) * limit)
	opts := options.Find().
		SetSort(bson.M{"createdAt": -1}).
		SetSkip(skip).
		SetLimit(int64(limit))

	c := database.DB.Collection("reviews")
	total, err := c.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	cursor, err := c.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reviews []models.Review
	if err := cursor.All(ctx, &reviews); err != nil {
		return nil, err
	}
	if reviews == nil {
		reviews = []models.Review{}
	}

	return &ReviewsResult{
		Reviews:      reviews,
		TotalPages:   int(math.Ceil(float64(total) / float64(limit))),
		CurrentPage:  page,
		TotalReviews: total,
	}, nil
}

func CreateReview(rv *models.Review) (*models.Review, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rv.Status = "Pending"
	rv.CreatedAt = time.Now()
	rv.UpdatedAt = time.Now()

	res, err := database.DB.Collection("reviews").InsertOne(ctx, rv)
	if err != nil {
		return nil, err
	}
	rv.ID = res.InsertedID.(bson.ObjectID)
	return rv, nil
}

func UpdateReviewStatus(id, status string) (*models.Review, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var review models.Review
	err = database.DB.Collection("reviews").FindOneAndUpdate(
		ctx,
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{"status": status, "updatedAt": time.Now()}},
		opts,
	).Decode(&review)
	return &review, err
}
