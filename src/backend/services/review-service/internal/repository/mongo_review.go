package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rentalflow/review-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoReviewRepository struct {
	coll *mongo.Collection
}

func NewMongoReviewRepository(db *mongo.Database) *MongoReviewRepository {
	return &MongoReviewRepository{
		coll: db.Collection("reviews"),
	}
}

func (r *MongoReviewRepository) Create(ctx context.Context, review *domain.Review) error {
	_, err := r.coll.InsertOne(ctx, review)
	return err
}

func (r *MongoReviewRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Review, error) {
	var review domain.Review
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&review)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrReviewNotFound
		}
		return nil, err
	}
	return &review, nil
}

func (r *MongoReviewRepository) GetByItem(ctx context.Context, itemID uuid.UUID, offset, limit int) ([]*domain.Review, int, error) {
	filter := bson.M{
		"target_item_id": itemID,
		"is_visible":     true,
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var reviews []*domain.Review
	if err := cursor.All(ctx, &reviews); err != nil {
		return nil, 0, err
	}
	return reviews, int(total), nil
}

func (r *MongoReviewRepository) GetByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*domain.Review, int, error) {
	filter := bson.M{
		"target_user_id": userID,
		"is_visible":     true,
	}

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var reviews []*domain.Review
	if err := cursor.All(ctx, &reviews); err != nil {
		return nil, 0, err
	}
	return reviews, int(total), nil
}

func (r *MongoReviewRepository) Update(ctx context.Context, review *domain.Review) error {
	update := bson.M{
		"$set": bson.M{
			"rating":     review.Rating,
			"comment":    review.Comment,
			"is_visible": review.IsVisible,
			"updated_at": time.Now(),
		},
	}
	result, err := r.coll.UpdateOne(ctx, bson.M{"_id": review.ID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return domain.ErrReviewNotFound
	}
	return nil
}

func (r *MongoReviewRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return domain.ErrReviewNotFound
	}
	return nil
}
