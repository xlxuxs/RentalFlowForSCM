package domain

import "errors"

// Domain errors
var (
	// Item errors
	ErrItemNotFound    = errors.New("rental item not found")
	ErrUnauthorized    = errors.New("unauthorized to perform this action")
	ErrInvalidCategory = errors.New("invalid item category")
	ErrInvalidPrice    = errors.New("invalid pricing information")

	// Availability errors
	ErrSlotNotFound     = errors.New("availability slot not found")
	ErrDateConflict     = errors.New("date range conflicts with existing bookings")
	ErrInvalidDateRange = errors.New("invalid date range")

	// Maintenance errors
	ErrMaintenanceNotFound = errors.New("maintenance log not found")
	ErrInvalidStatus       = errors.New("invalid status")
)
