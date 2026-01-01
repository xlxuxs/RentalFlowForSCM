package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rentalflow/inventory-service/internal/domain"
	"github.com/rentalflow/inventory-service/internal/repository"
)

// InventoryService handles inventory business logic
type InventoryService struct {
	itemRepo         repository.ItemRepository
	availabilityRepo repository.AvailabilityRepository
	maintenanceRepo  repository.MaintenanceRepository
}

// NewInventoryService creates a new inventory service
func NewInventoryService(
	itemRepo repository.ItemRepository,
	availabilityRepo repository.AvailabilityRepository,
	maintenanceRepo repository.MaintenanceRepository,
) *InventoryService {
	return &InventoryService{
		itemRepo:         itemRepo,
		availabilityRepo: availabilityRepo,
		maintenanceRepo:  maintenanceRepo,
	}
}

// CreateItem creates a new rental item
func (s *InventoryService) CreateItem(ctx context.Context, ownerID uuid.UUID, title, description string,
	category domain.ItemCategory, subcategory string, dailyRate, weeklyRate, monthlyRate, securityDeposit float64,
	location domain.Location, specs map[string]string, images []string) (*domain.RentalItem, error) {

	if !category.IsValid() {
		return nil, domain.ErrInvalidCategory
	}

	item := domain.NewRentalItem(ownerID, title, description, category, subcategory)
	item.DailyRate = dailyRate
	item.WeeklyRate = weeklyRate
	item.MonthlyRate = monthlyRate
	item.SecurityDeposit = securityDeposit
	item.Address = location.Address
	item.City = location.City
	item.Latitude = location.Latitude
	item.Longitude = location.Longitude
	item.Specifications = specs
	item.Images = images

	if err := s.itemRepo.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

// GetItem retrieves an item by ID
func (s *InventoryService) GetItem(ctx context.Context, itemID uuid.UUID) (*domain.RentalItem, error) {
	return s.itemRepo.GetByID(ctx, itemID)
}

// ListItems lists items with filters
func (s *InventoryService) ListItems(ctx context.Context, page, pageSize int, filters repository.ItemFilters) ([]*domain.RentalItem, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.itemRepo.List(ctx, offset, pageSize, filters)
}

// SearchItems searches items by query and filters
func (s *InventoryService) SearchItems(ctx context.Context, query string, page, pageSize int, filters repository.ItemFilters) ([]*domain.RentalItem, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.itemRepo.Search(ctx, query, filters, offset, pageSize)
}

// GetOwnerItems retrieves items owned by a specific user
func (s *InventoryService) GetOwnerItems(ctx context.Context, ownerID uuid.UUID, page, pageSize int) ([]*domain.RentalItem, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.itemRepo.GetByOwner(ctx, ownerID, offset, pageSize)
}

// UpdateItem updates an existing item
func (s *InventoryService) UpdateItem(ctx context.Context, itemID, ownerID uuid.UUID, updates map[string]interface{}) (*domain.RentalItem, error) {
	item, err := s.itemRepo.GetByID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	if item.OwnerID != ownerID {
		return nil, domain.ErrUnauthorized
	}

	// Apply updates
	if v, ok := updates["title"].(string); ok {
		item.Title = v
	}
	if v, ok := updates["description"].(string); ok {
		item.Description = v
	}
	if v, ok := updates["daily_rate"].(float64); ok {
		item.DailyRate = v
	}
	if v, ok := updates["is_active"].(bool); ok {
		item.IsActive = v
	}

	if err := s.itemRepo.Update(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

// DeleteItem deletes an item
func (s *InventoryService) DeleteItem(ctx context.Context, itemID, ownerID uuid.UUID) error {
	item, err := s.itemRepo.GetByID(ctx, itemID)
	if err != nil {
		return err
	}

	if item.OwnerID != ownerID {
		return domain.ErrUnauthorized
	}

	return s.itemRepo.Delete(ctx, itemID)
}

// BlockDates blocks dates for booking
func (s *InventoryService) BlockDates(ctx context.Context, itemID uuid.UUID, startDate, endDate time.Time, bookingID uuid.UUID) (*domain.AvailabilitySlot, error) {
	// Check for conflicts
	hasConflict, err := s.availabilityRepo.CheckConflict(ctx, itemID, startDate, endDate, nil)
	if err != nil {
		return nil, err
	}
	if hasConflict {
		return nil, domain.ErrDateConflict
	}

	slot := domain.NewAvailabilitySlot(itemID, startDate, endDate, domain.StatusBooked)
	slot.BookingID = &bookingID

	if err := s.availabilityRepo.Create(ctx, slot); err != nil {
		return nil, err
	}

	return slot, nil
}

// CreateMaintenanceLog creates a maintenance log
func (s *InventoryService) CreateMaintenanceLog(ctx context.Context, itemID, ownerID uuid.UUID, maintenanceType, description string, startDate time.Time, cost float64) (*domain.MaintenanceLog, error) {
	// Verify owner
	item, err := s.itemRepo.GetByID(ctx, itemID)
	if err != nil {
		return nil, err
	}
	if item.OwnerID != ownerID {
		return nil, domain.ErrUnauthorized
	}

	log := domain.NewMaintenanceLog(itemID, maintenanceType, description, startDate, cost)

	if err := s.maintenanceRepo.Create(ctx, log); err != nil {
		return nil, err
	}

	return log, nil
}
