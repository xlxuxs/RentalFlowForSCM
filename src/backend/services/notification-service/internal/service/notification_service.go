package service

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/rentalflow/notification-service/internal/domain"
	"github.com/rentalflow/notification-service/internal/repository"
)

type NotificationService struct {
	notificationRepo repository.NotificationRepository
	messageRepo      repository.MessageRepository
}

func NewNotificationService(notificationRepo repository.NotificationRepository, messageRepo repository.MessageRepository) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
		messageRepo:      messageRepo,
	}
}

func (s *NotificationService) HandleBookingEvent(ctx context.Context, eventData []byte) error {
	var event struct {
		ID           uuid.UUID `json:"id"`
		RenterID     uuid.UUID `json:"renter_id"`
		OwnerID      uuid.UUID `json:"owner_id"`
		RentalItemID uuid.UUID `json:"rental_item_id"`
		Status       string    `json:"status"`
	}

	if err := json.Unmarshal(eventData, &event); err != nil {
		return err
	}

	title := "Booking Update"
	message := ""
	var targetUserID uuid.UUID

	switch event.Status {
	case "pending": // booking.created
		targetUserID = event.OwnerID
		title = "New Booking Request"
		message = "You have a new booking request for your item."
	case "confirmed": // booking.confirmed
		targetUserID = event.RenterID
		title = "Booking Confirmed"
		message = "Your booking request has been confirmed by the owner."
	case "cancelled": // booking.cancelled
		targetUserID = event.RenterID
		title = "Booking Cancelled"
		message = "Your booking has been cancelled."
	}

	if targetUserID != uuid.Nil {
		_, err := s.SendNotification(ctx, targetUserID, "booking", title, message, domain.ChannelInApp)
		return err
	}

	return nil
}

func (s *NotificationService) SendNotification(ctx context.Context, userID uuid.UUID, notifType, title, message string, channel domain.NotificationChannel) (*domain.Notification, error) {
	notification := domain.NewNotification(userID, notifType, title, message, channel)
	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return nil, err
	}
	return notification, nil
}

func (s *NotificationService) GetUserNotifications(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*domain.Notification, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.notificationRepo.GetByUser(ctx, userID, offset, pageSize)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID uuid.UUID) error {
	return s.notificationRepo.MarkAsRead(ctx, notificationID)
}

func (s *NotificationService) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	return s.notificationRepo.GetUnreadCount(ctx, userID)
}

func (s *NotificationService) SendMessage(ctx context.Context, bookingID, senderID, receiverID uuid.UUID, content string) (*domain.Message, error) {
	message := domain.NewMessage(bookingID, senderID, receiverID, content)
	if err := s.messageRepo.Create(ctx, message); err != nil {
		return nil, err
	}
	return message, nil
}

func (s *NotificationService) GetBookingMessages(ctx context.Context, bookingID uuid.UUID, page, pageSize int) ([]*domain.Message, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}
	offset := (page - 1) * pageSize
	return s.messageRepo.GetByBooking(ctx, bookingID, offset, pageSize)
}
