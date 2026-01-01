package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rentalflow/payment-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoPaymentRepository struct {
	coll *mongo.Collection
}

func NewMongoPaymentRepository(db *mongo.Database) *MongoPaymentRepository {
	return &MongoPaymentRepository{
		coll: db.Collection("payments"),
	}
}

func (r *MongoPaymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	_, err := r.coll.InsertOne(ctx, payment)
	return err
}

func (r *MongoPaymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&payment)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrPaymentNotFound
		}
		return nil, err
	}
	return &payment, nil
}

func (r *MongoPaymentRepository) GetByBooking(ctx context.Context, bookingID uuid.UUID) ([]*domain.Payment, error) {
	opts := options.Find().SetSort(bson.M{"created_at": -1})
	cursor, err := r.coll.Find(ctx, bson.M{"booking_id": bookingID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var payments []*domain.Payment
	if err := cursor.All(ctx, &payments); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *MongoPaymentRepository) Update(ctx context.Context, payment *domain.Payment) error {
	update := bson.M{
		"$set": bson.M{
			"status":                  payment.Status,
			"provider_transaction_id": payment.ProviderTransactionID,
			"updated_at":              time.Now(),
		},
	}
	result, err := r.coll.UpdateOne(ctx, bson.M{"_id": payment.ID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return domain.ErrPaymentNotFound
	}
	return nil
}
