package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rentalflow/notification-service/internal/domain"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *domain.Notification) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error)
	GetByUser(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*domain.Notification, int, error)
	MarkAsRead(ctx context.Context, id uuid.UUID) error
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
}

type MessageRepository interface {
	Create(ctx context.Context, message *domain.Message) error
	GetByBooking(ctx context.Context, bookingID uuid.UUID, offset, limit int) ([]*domain.Message, int, error)
}
