package domain

import (
	"time"

	"github.com/google/uuid"
)

type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusConfirmed BookingStatus = "confirmed"
	StatusActive    BookingStatus = "active"
	StatusCompleted BookingStatus = "completed"
	StatusCancelled BookingStatus = "cancelled"
)

type CancellationPolicy string

const (
	PolicyFlexible CancellationPolicy = "flexible"
	PolicyModerate CancellationPolicy = "moderate"
	PolicyStrict   CancellationPolicy = "strict"
)

type Booking struct {
	ID                 uuid.UUID          `json:"id" bson:"_id"`
	BookingNumber      string             `json:"booking_number" bson:"booking_number"`
	RenterID           uuid.UUID          `json:"renter_id" bson:"renter_id"`
	OwnerID            uuid.UUID          `json:"owner_id" bson:"owner_id"`
	RentalItemID       uuid.UUID          `json:"rental_item_id" bson:"rental_item_id"`
	Status             BookingStatus      `json:"status" bson:"status"`
	StartDate          time.Time          `json:"start_date" bson:"start_date"`
	EndDate            time.Time          `json:"end_date" bson:"end_date"`
	TotalDays          int                `json:"total_days" bson:"total_days"`
	DailyRate          float64            `json:"daily_rate" bson:"daily_rate"`
	Subtotal           float64            `json:"subtotal" bson:"subtotal"`
	SecurityDeposit    float64            `json:"security_deposit" bson:"security_deposit"`
	ServiceFee         float64            `json:"service_fee" bson:"service_fee"`
	TotalAmount        float64            `json:"total_amount" bson:"total_amount"`
	PickupAddress      string             `json:"pickup_address,omitempty" bson:"pickup_address,omitempty"`
	PickupNotes        string             `json:"pickup_notes,omitempty" bson:"pickup_notes,omitempty"`
	PickupTime         *time.Time         `json:"pickup_time,omitempty" bson:"pickup_time,omitempty"`
	ReturnAddress      string             `json:"return_address,omitempty" bson:"return_address,omitempty"`
	ReturnNotes        string             `json:"return_notes,omitempty" bson:"return_notes,omitempty"`
	ReturnTime         *time.Time         `json:"return_time,omitempty" bson:"return_time,omitempty"`
	CancellationPolicy CancellationPolicy `json:"cancellation_policy" bson:"cancellation_policy"`
	AgreementSigned    bool               `json:"agreement_signed" bson:"agreement_signed"`
	AgreementURL       string             `json:"agreement_url,omitempty" bson:"agreement_url,omitempty"`
	CancelledBy        *uuid.UUID         `json:"cancelled_by,omitempty" bson:"cancelled_by,omitempty"`
	CancellationReason string             `json:"cancellation_reason,omitempty" bson:"cancellation_reason,omitempty"`
	PaymentStatus      string             `json:"payment_status,omitempty" bson:"payment_status,omitempty"`
	PaymentID          *uuid.UUID         `json:"payment_id,omitempty" bson:"payment_id,omitempty"`
	CreatedAt          time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at" bson:"updated_at"`
}

func NewBooking(renterID, ownerID, rentalItemID uuid.UUID, startDate, endDate time.Time, dailyRate, securityDeposit float64) *Booking {
	totalDays := int(endDate.Sub(startDate).Hours() / 24)
	if totalDays < 1 {
		totalDays = 1
	}

	subtotal := float64(totalDays) * dailyRate
	serviceFee := subtotal * 0.10 // 10%
	totalAmount := subtotal + serviceFee + securityDeposit

	now := time.Now()
	return &Booking{
		ID:                 uuid.New(),
		BookingNumber:      generateBookingNumber(),
		RenterID:           renterID,
		OwnerID:            ownerID,
		RentalItemID:       rentalItemID,
		Status:             StatusPending,
		StartDate:          startDate,
		EndDate:            endDate,
		TotalDays:          totalDays,
		DailyRate:          dailyRate,
		Subtotal:           subtotal,
		SecurityDeposit:    securityDeposit,
		ServiceFee:         serviceFee,
		TotalAmount:        totalAmount,
		CancellationPolicy: PolicyModerate,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
}

func generateBookingNumber() string {
	return "BK" + time.Now().Format("20060102") + uuid.New().String()[:4]
}
