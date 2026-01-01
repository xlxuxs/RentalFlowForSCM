package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rentalflow/booking-service/internal/domain"
	"github.com/rentalflow/booking-service/internal/repository"
	"github.com/rentalflow/rentalflow/pkg/messaging"
)

type BookingService struct {
	bookingRepo repository.BookingRepository
	broker      *messaging.MessageBroker
}

func NewBookingService(bookingRepo repository.BookingRepository, broker *messaging.MessageBroker) *BookingService {
	return &BookingService{
		bookingRepo: bookingRepo,
		broker:      broker,
	}
}

func (s *BookingService) CreateBooking(ctx context.Context, renterID, ownerID, rentalItemID uuid.UUID, startDate, endDate time.Time, dailyRate, securityDeposit float64) (*domain.Booking, error) {
	if endDate.Before(startDate) {
		return nil, domain.ErrInvalidDates
	}

	booking := domain.NewBooking(renterID, ownerID, rentalItemID, startDate, endDate, dailyRate, securityDeposit)
	if err := s.bookingRepo.Create(ctx, booking); err != nil {
		return nil, err
	}

	// Publish event
	if s.broker != nil {
		s.broker.Publish(ctx, "booking_events", "booking.created", booking)
	}

	return booking, nil
}

func (s *BookingService) GetBooking(ctx context.Context, bookingID uuid.UUID) (*domain.Booking, error) {
	return s.bookingRepo.GetByID(ctx, bookingID)
}

func (s *BookingService) GetRenterBookings(ctx context.Context, renterID uuid.UUID, page, pageSize int) ([]*domain.Booking, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.bookingRepo.GetByRenter(ctx, renterID, offset, pageSize)
}

func (s *BookingService) GetOwnerBookings(ctx context.Context, ownerID uuid.UUID, page, pageSize int) ([]*domain.Booking, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.bookingRepo.GetByOwner(ctx, ownerID, offset, pageSize)
}

func (s *BookingService) ConfirmBooking(ctx context.Context, bookingID, ownerID uuid.UUID) (*domain.Booking, error) {
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	if booking.OwnerID != ownerID {
		return nil, domain.ErrUnauthorized
	}

	if booking.Status != domain.StatusPending {
		return nil, domain.ErrInvalidStatus
	}

	booking.Status = domain.StatusConfirmed
	if err := s.bookingRepo.Update(ctx, booking); err != nil {
		return nil, err
	}

	// Publish event
	if s.broker != nil {
		s.broker.Publish(ctx, "booking_events", "booking.confirmed", booking)
	}

	return booking, nil
}

func (s *BookingService) CancelBooking(ctx context.Context, bookingID, userID uuid.UUID, reason string) (*domain.Booking, error) {
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	if booking.RenterID != userID && booking.OwnerID != userID {
		return nil, domain.ErrUnauthorized
	}

	if booking.Status == domain.StatusCancelled {
		return nil, domain.ErrAlreadyCancelled
	}

	if booking.Status == domain.StatusCompleted {
		return nil, domain.ErrCannotCancel
	}

	booking.Status = domain.StatusCancelled
	booking.CancelledBy = &userID
	booking.CancellationReason = reason
	if err := s.bookingRepo.Update(ctx, booking); err != nil {
		return nil, err
	}

	// Publish event
	if s.broker != nil {
		s.broker.Publish(ctx, "booking_events", "booking.cancelled", booking)
	}

	return booking, nil
}
