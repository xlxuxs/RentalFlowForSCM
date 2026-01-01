package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/rentalflow/booking-service/internal/domain"
	"github.com/rentalflow/booking-service/internal/service"
)

type HTTPHandler struct {
	bookingService *service.BookingService
}

func NewHTTPHandler(bookingService *service.BookingService) *HTTPHandler {
	return &HTTPHandler{bookingService: bookingService}
}

func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/ready", h.Ready)
	mux.HandleFunc("/api/bookings", h.HandleBookings)
	mux.HandleFunc("/api/bookings/renter", h.GetRenterBookings)
	mux.HandleFunc("/api/bookings/owner", h.GetOwnerBookings)
	mux.HandleFunc("/api/bookings/confirm", h.ConfirmBooking)
	mux.HandleFunc("/api/bookings/cancel", h.CancelBooking)
}

func (h *HTTPHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func (h *HTTPHandler) Ready(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}

func (h *HTTPHandler) HandleBookings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateBooking(w, r)
	case http.MethodGet:
		bookingID := r.URL.Query().Get("id")
		if bookingID != "" {
			h.GetBooking(w, r, bookingID)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "booking_id required"})
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *HTTPHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RenterID        string  `json:"renter_id"`
		OwnerID         string  `json:"owner_id"`
		RentalItemID    string  `json:"rental_item_id"`
		StartDate       string  `json:"start_date"`
		EndDate         string  `json:"end_date"`
		DailyRate       float64 `json:"daily_rate"`
		SecurityDeposit float64 `json:"security_deposit"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	renterID, _ := uuid.Parse(req.RenterID)
	ownerID, _ := uuid.Parse(req.OwnerID)
	rentalItemID, _ := uuid.Parse(req.RentalItemID)
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	booking, err := h.bookingService.CreateBooking(r.Context(), renterID, ownerID, rentalItemID, startDate, endDate, req.DailyRate, req.SecurityDeposit)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":             booking.ID.String(),
		"booking_number": booking.BookingNumber,
		"status":         booking.Status,
		"total_amount":   booking.TotalAmount,
		"start_date":     booking.StartDate,
		"end_date":       booking.EndDate,
	})
}

func (h *HTTPHandler) GetBooking(w http.ResponseWriter, r *http.Request, bookingID string) {
	id, err := uuid.Parse(bookingID)
	if err != nil {
		http.Error(w, "Invalid booking_id", http.StatusBadRequest)
		return
	}

	booking, err := h.bookingService.GetBooking(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":               booking.ID.String(),
		"booking_number":   booking.BookingNumber,
		"renter_id":        booking.RenterID.String(),
		"owner_id":         booking.OwnerID.String(),
		"rental_item_id":   booking.RentalItemID.String(),
		"status":           booking.Status,
		"start_date":       booking.StartDate,
		"end_date":         booking.EndDate,
		"total_days":       booking.TotalDays,
		"daily_rate":       booking.DailyRate,
		"total_amount":     booking.TotalAmount,
		"agreement_signed": booking.AgreementSigned,
	})
}

func (h *HTTPHandler) GetRenterBookings(w http.ResponseWriter, r *http.Request) {
	renterID := r.URL.Query().Get("renter_id")
	rid, _ := uuid.Parse(renterID)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	bookings, total, err := h.bookingService.GetRenterBookings(r.Context(), rid, page, pageSize)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"bookings": bookings,
		"total":    total,
	})
}

func (h *HTTPHandler) GetOwnerBookings(w http.ResponseWriter, r *http.Request) {
	ownerID := r.URL.Query().Get("owner_id")
	oid, _ := uuid.Parse(ownerID)
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	bookings, total, err := h.bookingService.GetOwnerBookings(r.Context(), oid, page, pageSize)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"bookings": bookings,
		"total":    total,
	})
}

func (h *HTTPHandler) ConfirmBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		BookingID string `json:"booking_id"`
		OwnerID   string `json:"owner_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bookingID, _ := uuid.Parse(req.BookingID)
	ownerID, _ := uuid.Parse(req.OwnerID)

	booking, err := h.bookingService.ConfirmBooking(r.Context(), bookingID, ownerID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     booking.ID.String(),
		"status": booking.Status,
	})
}

func (h *HTTPHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		BookingID string `json:"booking_id"`
		UserID    string `json:"user_id"`
		Reason    string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bookingID, _ := uuid.Parse(req.BookingID)
	userID, _ := uuid.Parse(req.UserID)

	booking, err := h.bookingService.CancelBooking(r.Context(), bookingID, userID, req.Reason)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     booking.ID.String(),
		"status": booking.Status,
	})
}

func (h *HTTPHandler) handleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch err {
	case domain.ErrBookingNotFound:
		w.WriteHeader(http.StatusNotFound)
	case domain.ErrUnauthorized:
		w.WriteHeader(http.StatusForbidden)
	case domain.ErrInvalidStatus, domain.ErrInvalidDates:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
