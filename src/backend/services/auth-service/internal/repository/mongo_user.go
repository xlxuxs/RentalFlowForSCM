package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rentalflow/auth-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoUserRepository implements UserRepository using MongoDB
type MongoUserRepository struct {
	coll *mongo.Collection
}

// NewMongoUserRepository creates a new MongoDB user repository
func NewMongoUserRepository(db *mongo.Database) *MongoUserRepository {
	return &MongoUserRepository{
		coll: db.Collection("users"),
	}
}

// Create creates a new user
func (r *MongoUserRepository) Create(ctx context.Context, user *domain.User) error {
	_, err := r.coll.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *MongoUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *MongoUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *MongoUserRepository) Update(ctx context.Context, user *domain.User) error {
	update := bson.M{
		"$set": bson.M{
			"email":                    user.Email,
			"password_hash":            user.PasswordHash,
			"first_name":               user.FirstName,
			"last_name":                user.LastName,
			"phone":                    user.Phone,
			"role":                     user.Role,
			"identity_verified":        user.IdentityVerified,
			"verification_status":      user.VerificationStatus,
			"refresh_token_hash":       user.RefreshTokenHash,
			"refresh_token_expires_at": user.RefreshTokenExpiresAt,
			"updated_at":               time.Now(),
		},
	}

	result, err := r.coll.UpdateOne(ctx, bson.M{"_id": user.ID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// Delete deletes a user
func (r *MongoUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// List retrieves a paginated list of users
func (r *MongoUserRepository) List(ctx context.Context, offset, limit int, filters UserFilters) ([]*domain.User, int, error) {
	filter := bson.M{}

	if filters.Role != nil {
		filter["role"] = *filters.Role
	}

	if filters.VerificationStatus != nil {
		filter["verification_status"] = *filters.VerificationStatus
	}

	// Get total count
	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Get users
	opts := options.Find().
		SetSort(bson.M{"created_at": -1}).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []*domain.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	return users, int(total), nil
}

// UpdateRefreshToken updates the refresh token for a user
func (r *MongoUserRepository) UpdateRefreshToken(ctx context.Context, userID uuid.UUID, hash string, expiresAt *time.Time) error {
	update := bson.M{
		"$set": bson.M{
			"refresh_token_hash":       hash,
			"refresh_token_expires_at": expiresAt,
			"updated_at":               time.Now(),
		},
	}

	result, err := r.coll.UpdateOne(ctx, bson.M{"_id": userID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// ClearRefreshToken clears the refresh token for a user
func (r *MongoUserRepository) ClearRefreshToken(ctx context.Context, userID uuid.UUID) error {
	return r.UpdateRefreshToken(ctx, userID, "", nil)
}
