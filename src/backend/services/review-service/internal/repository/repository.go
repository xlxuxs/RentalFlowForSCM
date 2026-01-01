package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rentalflow/review-service/internal/domain"
)

type ReviewRepository interface {
	Create(ctx context.Context, review *domain.Review) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Review, error)
	GetByItem(ctx context.Context, itemID uuid.UUID, offset, limit int) ([]*domain.Review, int, error)
	GetByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*domain.Review, int, error)
	Update(ctx context.Context, review *domain.Review) error
	Delete(ctx context.Context, id uuid.UUID) error
}
