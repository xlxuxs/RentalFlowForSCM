package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rentalflow/notification-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoNotificationRepository struct {
	coll *mongo.Collection
}

func NewMongoNotificationRepository(db *mongo.Database) *MongoNotificationRepository {
	return &MongoNotificationRepository{
		coll: db.Collection("notifications"),
	}
}

func (r *MongoNotificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	_, err := r.coll.InsertOne(ctx, notification)
	return err
}

func (r *MongoNotificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	var n domain.Notification
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&n)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNotificationNotFound
		}
		return nil, err
	}
	return &n, nil
}

func (r *MongoNotificationRepository) GetByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*domain.Notification, int, error) {
	filter := bson.M{"user_id": userID}

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

	var notifications []*domain.Notification
	if err := cursor.All(ctx, &notifications); err != nil {
		return nil, 0, err
	}
	return notifications, int(total), nil
}

func (r *MongoNotificationRepository) MarkAsRead(ctx context.Context, id uuid.UUID) error {
	update := bson.M{
		"$set": bson.M{
			"status":  domain.StatusRead,
			"read_at": time.Now(),
		},
	}
	result, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return domain.ErrNotificationNotFound
	}
	return nil
}

func (r *MongoNotificationRepository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	filter := bson.M{
		"user_id": userID,
		"status":  bson.M{"$ne": domain.StatusRead},
	}
	count, err := r.coll.CountDocuments(ctx, filter)
	return int(count), err
}

type MongoMessageRepository struct {
	coll *mongo.Collection
}

func NewMongoMessageRepository(db *mongo.Database) *MongoMessageRepository {
	return &MongoMessageRepository{
		coll: db.Collection("messages"),
	}
}

func (r *MongoMessageRepository) Create(ctx context.Context, message *domain.Message) error {
	_, err := r.coll.InsertOne(ctx, message)
	return err
}

func (r *MongoMessageRepository) GetByBooking(ctx context.Context, bookingID uuid.UUID, offset, limit int) ([]*domain.Message, int, error) {
	filter := bson.M{"booking_id": bookingID}

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

	var messages []*domain.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, 0, err
	}
	return messages, int(total), nil
}
