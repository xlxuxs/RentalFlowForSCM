package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rentalflow/auth-service/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDocumentRepository implements DocumentRepository using MongoDB
type MongoDocumentRepository struct {
	coll *mongo.Collection
}

// NewMongoDocumentRepository creates a new MongoDB document repository
func NewMongoDocumentRepository(db *mongo.Database) *MongoDocumentRepository {
	return &MongoDocumentRepository{
		coll: db.Collection("identity_documents"),
	}
}

// Create creates a new identity document
func (r *MongoDocumentRepository) Create(ctx context.Context, doc *domain.IdentityDocument) error {
	_, err := r.coll.InsertOne(ctx, doc)
	return err
}

// GetByID retrieves a document by ID
func (r *MongoDocumentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.IdentityDocument, error) {
	var doc domain.IdentityDocument
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("document not found")
		}
		return nil, err
	}
	return &doc, nil
}

// GetByUserID retrieves all documents for a user
func (r *MongoDocumentRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.IdentityDocument, error) {
	cursor, err := r.coll.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []*domain.IdentityDocument
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}
	return docs, nil
}

// Delete deletes a document
func (r *MongoDocumentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("document not found")
	}
	return nil
}
