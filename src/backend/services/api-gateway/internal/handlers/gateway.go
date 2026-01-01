package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rentalflow/api-gateway/config"
	"github.com/rentalflow/api-gateway/internal/clients"
)

type Gateway struct {
	authClient         *clients.HTTPClient
	inventoryClient    *clients.HTTPClient
	bookingClient      *clients.HTTPClient
	paymentClient      *clients.HTTPClient
	notificationClient *clients.HTTPClient
	reviewClient       *clients.HTTPClient
}

func NewGateway(cfg *config.Config) *Gateway {
	return &Gateway{
		authClient:         clients.NewHTTPClient(cfg.AuthServiceURL),
		inventoryClient:    clients.NewHTTPClient(cfg.InventoryServiceURL),
		bookingClient:      clients.NewHTTPClient(cfg.BookingServiceURL),
		paymentClient:      clients.NewHTTPClient(cfg.PaymentServiceURL),
		notificationClient: clients.NewHTTPClient(cfg.NotificationServiceURL),
		reviewClient:       clients.NewHTTPClient(cfg.ReviewServiceURL),
	}
}

func (g *Gateway) RegisterRoutes(r *mux.Router) {
	// Health check
	r.HandleFunc("/health", g.Health).Methods("GET")
	r.HandleFunc("/api/health", g.Health).Methods("GET")

	// Auth service routes (public)
	r.HandleFunc("/api/auth/register", g.forwardToAuth).Methods("POST")
	r.HandleFunc("/api/auth/login", g.forwardToAuth).Methods("POST")
	r.HandleFunc("/api/auth/validate", g.forwardToAuth).Methods("POST")
	r.HandleFunc("/api/auth/logout", g.forwardToAuth).Methods("POST")
	r.HandleFunc("/api/auth/profile", g.forwardToAuth).Methods("GET", "PUT")
	r.HandleFunc("/api/auth/avatar", g.forwardToAuth).Methods("POST")
	r.HandleFunc("/api/auth/change-password", g.forwardToAuth).Methods("POST")
	r.HandleFunc("/api/users", g.forwardToAuth).Methods("GET")

	// Inventory service routes
	r.PathPrefix("/api/inventory").HandlerFunc(g.forwardToInventory)
	r.PathPrefix("/api/items").HandlerFunc(g.forwardToInventory)

	// Booking service routes
	r.PathPrefix("/api/bookings").HandlerFunc(g.forwardToBooking)

	// Payment service routes
	r.PathPrefix("/api/payments").HandlerFunc(g.forwardToPayment)

	// Notification service routes
	r.PathPrefix("/api/notifications").HandlerFunc(g.forwardToNotification)
	r.PathPrefix("/api/messages").HandlerFunc(g.forwardToNotification)

	// Review service routes
	r.PathPrefix("/api/reviews").HandlerFunc(g.forwardToReview)
}

func (g *Gateway) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","service":"api-gateway"}`))
}

func (g *Gateway) forwardToAuth(w http.ResponseWriter, r *http.Request) {
	resp, err := g.authClient.Forward(r.Method, r.URL.Path, r.Body, getHeaders(r))
	if err != nil {
		clients.WriteError(w, http.StatusBadGateway, "Auth service unavailable")
		return
	}
	clients.ForwardResponse(w, resp)
}

func (g *Gateway) forwardToInventory(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if r.URL.RawQuery != "" {
		path += "?" + r.URL.RawQuery
	}
	resp, err := g.inventoryClient.Forward(r.Method, path, r.Body, getHeaders(r))
	if err != nil {
		clients.WriteError(w, http.StatusBadGateway, "Inventory service unavailable")
		return
	}
	clients.ForwardResponse(w, resp)
}

func (g *Gateway) forwardToBooking(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if r.URL.RawQuery != "" {
		path += "?" + r.URL.RawQuery
	}
	resp, err := g.bookingClient.Forward(r.Method, path, r.Body, getHeaders(r))
	if err != nil {
		clients.WriteError(w, http.StatusBadGateway, "Booking service unavailable")
		return
	}
	clients.ForwardResponse(w, resp)
}

func (g *Gateway) forwardToPayment(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if r.URL.RawQuery != "" {
		path += "?" + r.URL.RawQuery
	}
	resp, err := g.paymentClient.Forward(r.Method, path, r.Body, getHeaders(r))
	if err != nil {
		clients.WriteError(w, http.StatusBadGateway, "Payment service unavailable")
		return
	}
	clients.ForwardResponse(w, resp)
}

func (g *Gateway) forwardToNotification(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if r.URL.RawQuery != "" {
		path += "?" + r.URL.RawQuery
	}
	resp, err := g.notificationClient.Forward(r.Method, path, r.Body, getHeaders(r))
	if err != nil {
		clients.WriteError(w, http.StatusBadGateway, "Notification service unavailable")
		return
	}
	clients.ForwardResponse(w, resp)
}

func (g *Gateway) forwardToReview(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if r.URL.RawQuery != "" {
		path += "?" + r.URL.RawQuery
	}
	resp, err := g.reviewClient.Forward(r.Method, path, r.Body, getHeaders(r))
	if err != nil {
		clients.WriteError(w, http.StatusBadGateway, "Review service unavailable")
		return
	}
	clients.ForwardResponse(w, resp)
}

func getHeaders(r *http.Request) map[string]string {
	headers := make(map[string]string)
	for key, values := range r.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}
	return headers
}
