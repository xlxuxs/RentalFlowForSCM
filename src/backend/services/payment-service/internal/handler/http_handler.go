package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/rentalflow/payment-service/internal/domain"
	"github.com/rentalflow/payment-service/internal/service"
)

type HTTPHandler struct {
	paymentService *service.PaymentService
}

func NewHTTPHandler(paymentService *service.PaymentService) *HTTPHandler {
	return &HTTPHandler{paymentService: paymentService}
}

func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/ready", h.Ready)
	mux.HandleFunc("/api/payments/initialize", h.InitializePayment)
	mux.HandleFunc("/api/payments", h.GetPayment)
	mux.HandleFunc("/api/payments/booking", h.GetBookingPayments)
	mux.HandleFunc("/api/payments/refund", h.ProcessRefund)
	mux.HandleFunc("/api/payments/status", h.UpdateStatus)
	mux.HandleFunc("/api/payments/verify", h.VerifyPayment)
}

func (h *HTTPHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func (h *HTTPHandler) Ready(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}

func (h *HTTPHandler) InitializePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		BookingID string  `json:"booking_id"`
		UserID    string  `json:"user_id"`
		Amount    float64 `json:"amount"`
		Method    string  `json:"method"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bookingID, _ := uuid.Parse(req.BookingID)
	userID, _ := uuid.Parse(req.UserID)
	method := domain.PaymentMethod(req.Method)

	payment, err := h.paymentService.InitializePayment(r.Context(), bookingID, userID, req.Amount, method)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"payment_id":     payment.ID.String(),
		"checkout_url":   payment.CheckoutURL,
		"transaction_id": payment.ProviderTransactionID,
		"status":         payment.Status,
	})
}

func (h *HTTPHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	paymentID := r.URL.Query().Get("id")
	id, err := uuid.Parse(paymentID)
	if err != nil {
		http.Error(w, "Invalid payment_id", http.StatusBadRequest)
		return
	}

	payment, err := h.paymentService.GetPayment(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}

func (h *HTTPHandler) GetBookingPayments(w http.ResponseWriter, r *http.Request) {
	bookingID := r.URL.Query().Get("booking_id")
	bid, _ := uuid.Parse(bookingID)

	payments, err := h.paymentService.GetBookingPayments(r.Context(), bid)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"payments": payments,
		"count":    len(payments),
	})
}

func (h *HTTPHandler) ProcessRefund(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PaymentID string  `json:"payment_id"`
		Amount    float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	paymentID, _ := uuid.Parse(req.PaymentID)
	payment, err := h.paymentService.ProcessRefund(r.Context(), paymentID, req.Amount)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     payment.ID.String(),
		"status": payment.Status,
	})
}

func (h *HTTPHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		PaymentID     string `json:"payment_id"`
		Status        string `json:"status"`
		TransactionID string `json:"transaction_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	paymentID, _ := uuid.Parse(req.PaymentID)
	status := domain.PaymentStatus(req.Status)

	payment, err := h.paymentService.UpdatePaymentStatus(r.Context(), paymentID, status, req.TransactionID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     payment.ID.String(),
		"status": payment.Status,
	})
}

func (h *HTTPHandler) VerifyPayment(w http.ResponseWriter, r *http.Request) {
	txRef := r.URL.Query().Get("tx_ref")
	if txRef == "" {
		http.Error(w, "Missing tx_ref parameter", http.StatusBadRequest)
		return
	}

	chapaResp, err := h.paymentService.VerifyPayment(r.Context(), txRef)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tx_ref":    chapaResp.Data.TxRef,
		"reference": chapaResp.Data.Reference,
		"amount":    chapaResp.Data.Amount,
		"status":    chapaResp.Data.Status,
		"email":     chapaResp.Data.Email,
	})
}

func (h *HTTPHandler) handleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch err {
	case domain.ErrPaymentNotFound:
		w.WriteHeader(http.StatusNotFound)
	case domain.ErrInvalidAmount, domain.ErrInvalidPaymentMethod:
		w.WriteHeader(http.StatusBadRequest)
	case domain.ErrUnauthorized:
		w.WriteHeader(http.StatusForbidden)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
