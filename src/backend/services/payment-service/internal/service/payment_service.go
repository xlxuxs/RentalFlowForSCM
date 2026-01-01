package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rentalflow/payment-service/internal/chapa"
	"github.com/rentalflow/payment-service/internal/domain"
	"github.com/rentalflow/payment-service/internal/repository"
)

type PaymentService struct {
	paymentRepo repository.PaymentRepository
	chapaClient *chapa.Client
}

func NewPaymentService(paymentRepo repository.PaymentRepository, chapaClient *chapa.Client) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		chapaClient: chapaClient,
	}
}

func (s *PaymentService) InitializePayment(ctx context.Context, bookingID, userID uuid.UUID, amount float64, method domain.PaymentMethod) (*domain.Payment, error) {
	if amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	payment := domain.NewPayment(bookingID, userID, amount, method)
	payment.PaymentType = "booking"
	payment.ProviderName = string(method)

	// Handle Chapa payments
	if method == domain.MethodChapa && s.chapaClient != nil {
		txRef := fmt.Sprintf("RF-%s-%s", bookingID.String()[:8], payment.ID.String()[:8])

		chapaReq := chapa.InitializePaymentRequest{
			Amount:      amount,
			Currency:    "ETB",
			Email:       "customer@rentalflow.com", // TODO: Get from user profile
			FirstName:   "RentalFlow",
			LastName:    "Customer",
			TxRef:       txRef,
			CallbackURL: fmt.Sprintf("http://localhost:3001/payment/callback?tx_ref=%s", txRef),
			ReturnURL:   fmt.Sprintf("http://localhost:3001/payment/callback?tx_ref=%s&booking_id=%s", txRef, bookingID.String()),
			CustomTitle: "RentalFlow Payment",
			CustomDesc:  fmt.Sprintf("Payment for Booking #%s", bookingID.String()[:8]),
			Metadata: map[string]string{
				"booking_id": bookingID.String(),
				"payment_id": payment.ID.String(),
			},
		}

		chapaResp, err := s.chapaClient.InitializePayment(chapaReq)
		if err != nil {
			return nil, fmt.Errorf("chapa initialization failed: %w", err)
		}

		payment.CheckoutURL = chapaResp.Data.CheckoutURL
		payment.ProviderTransactionID = txRef
	} else {
		// Fallback for other payment methods
		payment.CheckoutURL = fmt.Sprintf("https://%s.com/checkout/%s", method, payment.ID.String())
		payment.ProviderTransactionID = fmt.Sprintf("TXN_%s", uuid.New().String()[:8])
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *PaymentService) GetPayment(ctx context.Context, paymentID uuid.UUID) (*domain.Payment, error) {
	return s.paymentRepo.GetByID(ctx, paymentID)
}

func (s *PaymentService) GetBookingPayments(ctx context.Context, bookingID uuid.UUID) ([]*domain.Payment, error) {
	return s.paymentRepo.GetByBooking(ctx, bookingID)
}

func (s *PaymentService) UpdatePaymentStatus(ctx context.Context, paymentID uuid.UUID, status domain.PaymentStatus, transactionID string) (*domain.Payment, error) {
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	payment.Status = status
	payment.ProviderTransactionID = transactionID

	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *PaymentService) VerifyPayment(ctx context.Context, txRef string) (*chapa.VerifyPaymentResponse, error) {
	if s.chapaClient == nil {
		return nil, fmt.Errorf("chapa client not initialized")
	}
	return s.chapaClient.VerifyPayment(txRef)
}

func (s *PaymentService) ProcessRefund(ctx context.Context, paymentID uuid.UUID, amount float64) (*domain.Payment, error) {
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, err
	}

	if payment.Status != domain.StatusCompleted {
		return nil, domain.ErrRefundNotAllowed
	}

	if amount > payment.Amount {
		return nil, domain.ErrInvalidAmount
	}

	payment.Status = domain.StatusRefunded
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		return nil, err
	}

	return payment, nil
}
