package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rentalflow/inventory-service/internal/domain"
)

// ItemRepository defines the interface for rental item data access
type ItemRepository interface {
	Create(ctx context.Context, item *domain.RentalItem) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.RentalItem, error)
	GetByOwner(ctx context.Context, ownerID uuid.UUID, offset, limit int) ([]*domain.RentalItem, int, error)
	List(ctx context.Context, offset, limit int, filters ItemFilters) ([]*domain.RentalItem, int, error)
	Search(ctx context.Context, query string, filters ItemFilters, offset, limit int) ([]*domain.RentalItem, int, error)
	Update(ctx context.Context, item *domain.RentalItem) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// ItemFilters defines filters for listing items
type ItemFilters struct {
	Category *domain.ItemCategory
	City     *string
	MinPrice *float64
	MaxPrice *float64
	IsActive *bool
	SortBy   *string
}

// AvailabilityRepository defines the interface for availability slot data access
type AvailabilityRepository interface {
	Create(ctx context.Context, slot *domain.AvailabilitySlot) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.AvailabilitySlot, error)
	GetByItem(ctx context.Context, itemID uuid.UUID, startDate, endDate time.Time) ([]*domain.AvailabilitySlot, error)
	Update(ctx context.Context, slot *domain.AvailabilitySlot) error
	Delete(ctx context.Context, id uuid.UUID) error
	CheckConflict(ctx context.Context, itemID uuid.UUID, startDate, endDate time.Time, excludeSlotID *uuid.UUID) (bool, error)
}

// MaintenanceRepository defines the interface for maintenance log data access
type MaintenanceRepository interface {
	Create(ctx context.Context, log *domain.MaintenanceLog) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.MaintenanceLog, error)
	GetByItem(ctx context.Context, itemID uuid.UUID, offset, limit int) ([]*domain.MaintenanceLog, int, error)
	Update(ctx context.Context, log *domain.MaintenanceLog) error
	Delete(ctx context.Context, id uuid.UUID) error
}
