package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rentalflow/auth-service/internal/domain"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *domain.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*domain.User, error)

	// Update updates a user
	Update(ctx context.Context, user *domain.User) error

	// Delete deletes a user
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves a paginated list of users
	List(ctx context.Context, offset, limit int, filters UserFilters) ([]*domain.User, int, error)

	// UpdateRefreshToken updates the refresh token for a user
	UpdateRefreshToken(ctx context.Context, userID uuid.UUID, hash string, expiresAt *time.Time) error

	// ClearRefreshToken clears the refresh token for a user (logout)
	ClearRefreshToken(ctx context.Context, userID uuid.UUID) error
}

// UserFilters defines filters for listing users
type UserFilters struct {
	Role               *domain.UserRole
	VerificationStatus *domain.VerificationStatus
}

// DocumentRepository defines the interface for identity document data access
type DocumentRepository interface {
	// Create creates a new identity document
	Create(ctx context.Context, doc *domain.IdentityDocument) error

	// GetByID retrieves a document by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.IdentityDocument, error)

	// GetByUserID retrieves all documents for a user
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.IdentityDocument, error)

	// Delete deletes a document
	Delete(ctx context.Context, id uuid.UUID) error
}
