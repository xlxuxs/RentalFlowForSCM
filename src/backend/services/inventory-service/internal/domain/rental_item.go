package domain

import (
	"time"

	"github.com/google/uuid"
)

// ItemCategory represents the category of a rental item
type ItemCategory string

const (
	CategoryVehicle   ItemCategory = "vehicle"
	CategoryEquipment ItemCategory = "equipment"
	CategoryProperty  ItemCategory = "property"
)

// IsValid checks if the category is valid
func (c ItemCategory) IsValid() bool {
	switch c {
	case CategoryVehicle, CategoryEquipment, CategoryProperty:
		return true
	}
	return false
}

// AvailabilityStatus represents the availability status
type AvailabilityStatus string

const (
	StatusAvailable   AvailabilityStatus = "available"
	StatusBooked      AvailabilityStatus = "booked"
	StatusMaintenance AvailabilityStatus = "maintenance"
	StatusBlocked     AvailabilityStatus = "blocked"
)

// MaintenanceStatus represents maintenance status
type MaintenanceStatus string

const (
	MaintenanceScheduled  MaintenanceStatus = "scheduled"
	MaintenanceInProgress MaintenanceStatus = "in_progress"
	MaintenanceCompleted  MaintenanceStatus = "completed"
)

// RentalItem represents a rental item
type RentalItem struct {
	ID          uuid.UUID    `json:"id" bson:"_id"`
	OwnerID     uuid.UUID    `json:"owner_id" bson:"owner_id"`
	Title       string       `json:"title" bson:"title"`
	Description string       `json:"description" bson:"description"`
	Category    ItemCategory `json:"category" bson:"category"`
	Subcategory string       `json:"subcategory" bson:"subcategory"`

	// Pricing
	DailyRate       float64 `json:"daily_rate" bson:"daily_rate"`
	WeeklyRate      float64 `json:"weekly_rate" bson:"weekly_rate"`
	MonthlyRate     float64 `json:"monthly_rate" bson:"monthly_rate"`
	SecurityDeposit float64 `json:"security_deposit" bson:"security_deposit"`

	// Location
	Address   string  `json:"address" bson:"address"`
	City      string  `json:"city" bson:"city"`
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`

	// Specifications (stored as map)
	Specifications map[string]string `json:"specifications" bson:"specifications"`

	// Images
	Images []string `json:"images" bson:"images"`

	IsActive  bool      `json:"is_active" bson:"is_active"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// NewRentalItem creates a new rental item
func NewRentalItem(ownerID uuid.UUID, title, description string, category ItemCategory, subcategory string) *RentalItem {
	now := time.Now()
	return &RentalItem{
		ID:             uuid.New(),
		OwnerID:        ownerID,
		Title:          title,
		Description:    description,
		Category:       category,
		Subcategory:    subcategory,
		Specifications: make(map[string]string),
		Images:         []string{},
		IsActive:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// AvailabilitySlot represents an availability slot for a rental item
type AvailabilitySlot struct {
	ID           uuid.UUID          `json:"id" bson:"_id"`
	RentalItemID uuid.UUID          `json:"rental_item_id" bson:"rental_item_id"`
	StartDate    time.Time          `json:"start_date" bson:"start_date"`
	EndDate      time.Time          `json:"end_date" bson:"end_date"`
	Status       AvailabilityStatus `json:"status" bson:"status"`
	BookingID    *uuid.UUID         `json:"booking_id,omitempty" bson:"booking_id,omitempty"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
}

// NewAvailabilitySlot creates a new availability slot
func NewAvailabilitySlot(rentalItemID uuid.UUID, startDate, endDate time.Time, status AvailabilityStatus) *AvailabilitySlot {
	return &AvailabilitySlot{
		ID:           uuid.New(),
		RentalItemID: rentalItemID,
		StartDate:    startDate,
		EndDate:      endDate,
		Status:       status,
		CreatedAt:    time.Now(),
	}
}

// MaintenanceLog represents a maintenance log entry
type MaintenanceLog struct {
	ID              uuid.UUID         `json:"id" bson:"_id"`
	RentalItemID    uuid.UUID         `json:"rental_item_id" bson:"rental_item_id"`
	MaintenanceType string            `json:"maintenance_type" bson:"maintenance_type"`
	Description     string            `json:"description" bson:"description"`
	StartDate       time.Time         `json:"start_date" bson:"start_date"`
	EndDate         *time.Time        `json:"end_date,omitempty" bson:"end_date,omitempty"`
	Cost            float64           `json:"cost" bson:"cost"`
	Status          MaintenanceStatus `json:"status" bson:"status"`
	CreatedAt       time.Time         `json:"created_at" bson:"created_at"`
}

// NewMaintenanceLog creates a new maintenance log
func NewMaintenanceLog(rentalItemID uuid.UUID, maintenanceType, description string, startDate time.Time, cost float64) *MaintenanceLog {
	return &MaintenanceLog{
		ID:              uuid.New(),
		RentalItemID:    rentalItemID,
		MaintenanceType: maintenanceType,
		Description:     description,
		StartDate:       startDate,
		Cost:            cost,
		Status:          MaintenanceScheduled,
		CreatedAt:       time.Now(),
	}
}

// Location represents a geographic location
type Location struct {
	Address   string
	City      string
	Latitude  float64
	Longitude float64
}
