package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rentalflow/booking-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoBookingRepository struct {
	coll *mongo.Collection
}

func NewMongoBookingRepository(db *mongo.Database) *MongoBookingRepository {
	return &MongoBookingRepository{
		coll: db.Collection("bookings"),
	}
}

func (r *MongoBookingRepository) Create(ctx context.Context, booking *domain.Booking) error {
	_, err := r.coll.InsertOne(ctx, booking)
	return err
}

func (r *MongoBookingRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	var booking domain.Booking
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&booking)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrBookingNotFound
		}
		return nil, err
	}
	return &booking, nil
}

func (r *MongoBookingRepository) GetByRenter(ctx context.Context, renterID uuid.UUID, offset, limit int) ([]*domain.Booking, int, error) {
	filter := bson.M{"renter_id": renterID}

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

	var bookings []*domain.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, 0, err
	}
	return bookings, int(total), nil
}

func (r *MongoBookingRepository) GetByOwner(ctx context.Context, ownerID uuid.UUID, offset, limit int) ([]*domain.Booking, int, error) {
	filter := bson.M{"owner_id": ownerID}

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

	var bookings []*domain.Booking
	if err := cursor.All(ctx, &bookings); err != nil {
		return nil, 0, err
	}
	return bookings, int(total), nil
}

func (r *MongoBookingRepository) Update(ctx context.Context, booking *domain.Booking) error {
	update := bson.M{
		"$set": bson.M{
			"status":           booking.Status,
			"agreement_signed": booking.AgreementSigned,
			"updated_at":       time.Now(),
			// Add other updatable fields as needed based on logic
		},
	}

	result, err := r.coll.UpdateOne(ctx, bson.M{"_id": booking.ID}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return domain.ErrBookingNotFound
	}
	return nil
}
