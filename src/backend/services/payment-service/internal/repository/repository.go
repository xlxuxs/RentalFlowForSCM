package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rentalflow/payment-service/internal/domain"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *domain.Payment) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error)
	GetByBooking(ctx context.Context, bookingID uuid.UUID) ([]*domain.Payment, error)
	Update(ctx context.Context, payment *domain.Payment) error
}
