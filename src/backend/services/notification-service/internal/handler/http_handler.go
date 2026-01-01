package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/rentalflow/notification-service/internal/email"
	"github.com/rentalflow/notification-service/internal/service"
)

type HTTPHandler struct {
	notificationService *service.NotificationService
	emailService        *email.Service
}

func NewHTTPHandler(notificationService *service.NotificationService, emailService *email.Service) *HTTPHandler {
	return &HTTPHandler{
		notificationService: notificationService,
		emailService:        emailService,
	}
}

func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	fmt.Println("Registering /health")
	mux.HandleFunc("/health", h.Health)
	fmt.Println("Registering /api/notifications/booking-created")
	mux.HandleFunc("/api/notifications/booking-created", h.SendBookingCreated)
	fmt.Println("Registering /api/notifications/payment-success", h.SendPaymentSuccess)
	mux.HandleFunc("/api/notifications/payment-success", h.SendPaymentSuccess)
	fmt.Println("Registering /api/notifications/review-received")
	mux.HandleFunc("/api/notifications/review-received", h.SendReviewReceived)

	// New In-App Notification Routes
	fmt.Println("Registering /api/notifications/user")
	mux.HandleFunc("/api/notifications/user", h.GetUserNotifications)
	fmt.Println("Registering /api/notifications/mark-read")
	mux.HandleFunc("/api/notifications/mark-read", h.MarkAsRead)
	fmt.Println("Registering /api/notifications/unread-count")
	mux.HandleFunc("/api/notifications/unread-count", h.GetUnreadCount)

	// Message Routes
	fmt.Println("Registering /api/messages/send")
	mux.HandleFunc("/api/messages/send", h.SendMessage)
	fmt.Println("Registering /api/messages/booking")
	mux.HandleFunc("/api/messages/booking", h.GetBookingMessages)
}

func (h *HTTPHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func (h *HTTPHandler) SendBookingCreated(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		To   string                 `json:"to"`
		Data map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.emailService.SendBookingCreated(req.To, req.Data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (h *HTTPHandler) SendPaymentSuccess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		To   string                 `json:"to"`
		Data map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.emailService.SendPaymentSuccess(req.To, req.Data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (h *HTTPHandler) SendReviewReceived(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		To   string                 `json:"to"`
		Data map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.emailService.SendReviewReceived(req.To, req.Data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
func (h *HTTPHandler) GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	notifications, total, err := h.notificationService.GetUserNotifications(r.Context(), userID, page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"notifications": notifications,
		"total":         total,
	})
}

func (h *HTTPHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		NotificationID string `json:"notification_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(req.NotificationID)
	if err != nil {
		http.Error(w, "Invalid notification_id", http.StatusBadRequest)
		return
	}

	if err := h.notificationService.MarkAsRead(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *HTTPHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	count, err := h.notificationService.GetUnreadCount(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"count": count})
}

func (h *HTTPHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		BookingID  string `json:"booking_id"`
		SenderID   string `json:"sender_id"`
		ReceiverID string `json:"receiver_id"`
		Content    string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	bid, _ := uuid.Parse(req.BookingID)
	sid, _ := uuid.Parse(req.SenderID)
	rid, _ := uuid.Parse(req.ReceiverID)

	message, err := h.notificationService.SendMessage(r.Context(), bid, sid, rid, req.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

func (h *HTTPHandler) GetBookingMessages(w http.ResponseWriter, r *http.Request) {
	bookingIDStr := r.URL.Query().Get("booking_id")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		http.Error(w, "Invalid booking_id", http.StatusBadRequest)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	messages, total, err := h.notificationService.GetBookingMessages(r.Context(), bookingID, page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"messages": messages,
		"total":    total,
	})
}
