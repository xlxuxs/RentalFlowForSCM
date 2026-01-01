package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rentalflow/booking-service/internal/domain"
)

type BookingRepository interface {
	Create(ctx context.Context, booking *domain.Booking) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error)
	GetByRenter(ctx context.Context, renterID uuid.UUID, offset, limit int) ([]*domain.Booking, int, error)
	GetByOwner(ctx context.Context, ownerID uuid.UUID, offset, limit int) ([]*domain.Booking, int, error)
	Update(ctx context.Context, booking *domain.Booking) error
}
