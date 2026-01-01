package domain

import (
	"time"

	"github.com/google/uuid"
)

type PaymentStatus string

const (
	StatusPending    PaymentStatus = "pending"
	StatusProcessing PaymentStatus = "processing"
	StatusCompleted  PaymentStatus = "completed"
	StatusFailed     PaymentStatus = "failed"
	StatusRefunded   PaymentStatus = "refunded"
)

type PaymentMethod string

const (
	MethodChapa        PaymentMethod = "chapa"
	MethodTelebirr     PaymentMethod = "telebirr"
	MethodBankTransfer PaymentMethod = "bank_transfer"
	MethodCash         PaymentMethod = "cash"
)

type Payment struct {
	ID                    uuid.UUID     `json:"id" bson:"_id"`
	BookingID             uuid.UUID     `json:"booking_id" bson:"booking_id"`
	UserID                uuid.UUID     `json:"user_id" bson:"user_id"`
	PaymentType           string        `json:"payment_type" bson:"payment_type"`
	Amount                float64       `json:"amount" bson:"amount"`
	Currency              string        `json:"currency" bson:"currency"`
	Status                PaymentStatus `json:"status" bson:"status"`
	Method                PaymentMethod `json:"method" bson:"method"`
	RentalFee             float64       `json:"rental_fee" bson:"rental_fee"`
	SecurityDeposit       float64       `json:"security_deposit" bson:"security_deposit"`
	ServiceFee            float64       `json:"service_fee" bson:"service_fee"`
	AdditionalServices    float64       `json:"additional_services" bson:"additional_services"`
	Tax                   float64       `json:"tax" bson:"tax"`
	DepositHeld           bool          `json:"deposit_held" bson:"deposit_held"`
	DepositStatus         string        `json:"deposit_status" bson:"deposit_status"`
	ProviderName          string        `json:"provider_name" bson:"provider_name"`
	ProviderTransactionID string        `json:"provider_transaction_id" bson:"provider_transaction_id"`
	CheckoutURL           string        `json:"checkout_url" bson:"checkout_url"`
	ReceiptURL            string        `json:"receipt_url" bson:"receipt_url"`
	CreatedAt             time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at" bson:"updated_at"`
}

func NewPayment(bookingID, userID uuid.UUID, amount float64, method PaymentMethod) *Payment {
	return &Payment{
		ID:        uuid.New(),
		BookingID: bookingID,
		UserID:    userID,
		Amount:    amount,
		Currency:  "ETB",
		Status:    StatusPending,
		Method:    method,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
