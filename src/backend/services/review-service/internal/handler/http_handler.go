package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/rentalflow/review-service/internal/domain"
	"github.com/rentalflow/review-service/internal/service"
)

type HTTPHandler struct {
	reviewService *service.ReviewService
}

func NewHTTPHandler(reviewService *service.ReviewService) *HTTPHandler {
	return &HTTPHandler{reviewService: reviewService}
}

func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/api/reviews", h.HandleReviews)
	mux.HandleFunc("/api/reviews/item", h.GetItemReviews)
	mux.HandleFunc("/api/reviews/user", h.GetUserReviews)
}

func (h *HTTPHandler) Health(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func (h *HTTPHandler) HandleReviews(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateReview(w, r)
	case http.MethodGet:
		reviewID := r.URL.Query().Get("id")
		if reviewID != "" {
			h.GetReview(w, r, reviewID)
		}
	case http.MethodPut:
		h.UpdateReview(w, r)
	case http.MethodDelete:
		h.DeleteReview(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *HTTPHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ItemID     string  `json:"item_id"`
		BookingID  string  `json:"booking_id"`
		ReviewerID string  `json:"reviewer_id"`
		ReviewType string  `json:"review_type"`
		Rating     float64 `json:"rating"`
		Comment    string  `json:"comment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	itemID, err := uuid.Parse(req.ItemID)
	if err != nil {
		http.Error(w, "Invalid item_id", http.StatusBadRequest)
		return
	}

	reviewerID, err := uuid.Parse(req.ReviewerID)
	if err != nil {
		http.Error(w, "Invalid reviewer_id", http.StatusBadRequest)
		return
	}

	var bookingIDPtr *uuid.UUID
	if req.BookingID != "" {
		bookingID, err := uuid.Parse(req.BookingID)
		if err == nil {
			bookingIDPtr = &bookingID
		}
	}

	// Default to renter_to_item if not provided
	rType := domain.TypeRenterToItem
	if req.ReviewType != "" {
		rType = domain.ReviewType(req.ReviewType)
	}

	review, err := h.reviewService.CreateReview(r.Context(), itemID, bookingIDPtr, reviewerID, rType, req.Rating, req.Comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
}

func (h *HTTPHandler) GetReview(w http.ResponseWriter, r *http.Request, reviewID string) {
	id, _ := uuid.Parse(reviewID)
	review, err := h.reviewService.GetReview(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

func (h *HTTPHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ReviewID string  `json:"review_id"`
		Rating   float64 `json:"rating"`
		Comment  string  `json:"comment"`
	}

	json.NewDecoder(r.Body).Decode(&req)
	reviewID, _ := uuid.Parse(req.ReviewID)

	review, err := h.reviewService.UpdateReview(r.Context(), reviewID, req.Rating, req.Comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

func (h *HTTPHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	reviewID := r.URL.Query().Get("id")
	id, _ := uuid.Parse(reviewID)

	if err := h.reviewService.DeleteReview(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

func (h *HTTPHandler) GetItemReviews(w http.ResponseWriter, r *http.Request) {
	itemID, _ := uuid.Parse(r.URL.Query().Get("item_id"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	reviews, total, err := h.reviewService.GetItemReviews(r.Context(), itemID, page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"reviews": reviews, "total": total})
}

func (h *HTTPHandler) GetUserReviews(w http.ResponseWriter, r *http.Request) {
	userID, _ := uuid.Parse(r.URL.Query().Get("user_id"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	reviews, total, err := h.reviewService.GetUserReviews(r.Context(), userID, page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"reviews": reviews, "total": total})
}
