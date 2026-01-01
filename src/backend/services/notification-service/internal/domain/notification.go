package domain

import (
	"time"

	"github.com/google/uuid"
)

type NotificationStatus string

const (
	StatusPending NotificationStatus = "pending"
	StatusSent    NotificationStatus = "sent"
	StatusFailed  NotificationStatus = "failed"
	StatusRead    NotificationStatus = "read"
)

type NotificationChannel string

const (
	ChannelEmail NotificationChannel = "email"
	ChannelSMS   NotificationChannel = "sms"
	ChannelPush  NotificationChannel = "push"
	ChannelInApp NotificationChannel = "in_app"
)

type Notification struct {
	ID               uuid.UUID           `json:"id" bson:"_id"`
	UserID           uuid.UUID           `json:"user_id" bson:"user_id"`
	NotificationType string              `json:"notification_type" bson:"notification_type"`
	Title            string              `json:"title" bson:"title"`
	Message          string              `json:"message" bson:"message"`
	Channel          NotificationChannel `json:"channel" bson:"channel"`
	Priority         string              `json:"priority" bson:"priority"`
	Status           NotificationStatus  `json:"status" bson:"status"`
	ActionURL        string              `json:"action_url" bson:"action_url"`
	SentAt           *time.Time          `json:"sent_at" bson:"sent_at"`
	ReadAt           *time.Time          `json:"read_at" bson:"read_at"`
	CreatedAt        time.Time           `json:"created_at" bson:"created_at"`
}

type Message struct {
	ID          uuid.UUID  `json:"id" bson:"_id"`
	BookingID   uuid.UUID  `json:"booking_id" bson:"booking_id"`
	SenderID    uuid.UUID  `json:"sender_id" bson:"sender_id"`
	ReceiverID  uuid.UUID  `json:"receiver_id" bson:"receiver_id"`
	Content     string     `json:"content" bson:"content"`
	Attachments []string   `json:"attachments" bson:"attachments"`
	IsRead      bool       `json:"is_read" bson:"is_read"`
	ReadAt      *time.Time `json:"read_at" bson:"read_at"`
	CreatedAt   time.Time  `json:"created_at" bson:"created_at"`
}

func NewNotification(userID uuid.UUID, notifType, title, message string, channel NotificationChannel) *Notification {
	return &Notification{
		ID:               uuid.New(),
		UserID:           userID,
		NotificationType: notifType,
		Title:            title,
		Message:          message,
		Channel:          channel,
		Priority:         "medium",
		Status:           StatusPending,
		CreatedAt:        time.Now(),
	}
}

func NewMessage(bookingID, senderID, receiverID uuid.UUID, content string) *Message {
	return &Message{
		ID:         uuid.New(),
		BookingID:  bookingID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		IsRead:     false,
		CreatedAt:  time.Now(),
	}
}
