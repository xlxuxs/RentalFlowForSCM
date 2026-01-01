package domain

import (
	"time"

	"github.com/google/uuid"
)

type ReviewType string

const (
	TypeRenterToOwner ReviewType = "renter_to_owner"
	TypeOwnerToRenter ReviewType = "owner_to_renter"
	TypeRenterToItem  ReviewType = "renter_to_item"
)

type Review struct {
	ID           uuid.UUID  `json:"id" bson:"_id"`
	BookingID    uuid.UUID  `json:"booking_id" bson:"booking_id"`
	ReviewerID   uuid.UUID  `json:"reviewer_id" bson:"reviewer_id"`
	TargetUserID *uuid.UUID `json:"target_user_id,omitempty" bson:"target_user_id,omitempty"`
	TargetItemID *uuid.UUID `json:"target_item_id,omitempty" bson:"target_item_id,omitempty"`
	ReviewType   ReviewType `json:"review_type" bson:"review_type"`
	Rating       float64    `json:"rating" bson:"rating"`
	Comment      string     `json:"comment" bson:"comment"`
	IsVerified   bool       `json:"is_verified" bson:"is_verified"`
	IsVisible    bool       `json:"is_visible" bson:"is_visible"`
	CreatedAt    time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" bson:"updated_at"`
}

func NewReview(bookingID, reviewerID uuid.UUID, reviewType ReviewType, rating float64, comment string) *Review {
	now := time.Now()
	if rating < 1.0 || rating > 5.0 {
		rating = 5.0
	}
	return &Review{
		ID:         uuid.New(),
		BookingID:  bookingID,
		ReviewerID: reviewerID,
		ReviewType: reviewType,
		Rating:     rating,
		Comment:    comment,
		IsVerified: false,
		IsVisible:  true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}
