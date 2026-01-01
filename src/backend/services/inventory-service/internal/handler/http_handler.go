package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/rentalflow/inventory-service/internal/domain"
	"github.com/rentalflow/inventory-service/internal/repository"
	"github.com/rentalflow/inventory-service/internal/service"
)

// HTTPHandler provides REST endpoints for testing
type HTTPHandler struct {
	inventoryService *service.InventoryService
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(inventoryService *service.InventoryService) *HTTPHandler {
	return &HTTPHandler{inventoryService: inventoryService}
}

// RegisterRoutes registers HTTP routes
func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/ready", h.Ready)
	mux.HandleFunc("/api/items", h.HandleItems)
	mux.HandleFunc("/api/items/owner", h.GetOwnerItems)
	mux.HandleFunc("/api/items/search", h.SearchItems)
	mux.HandleFunc("/api/availability/block", h.BlockDates)
	mux.HandleFunc("/api/maintenance", h.CreateMaintenance)
}

func (h *HTTPHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func (h *HTTPHandler) Ready(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}

type CreateItemRequest struct {
	OwnerID         string            `json:"owner_id"`
	Title           string            `json:"title"`
	Description     string            `json:"description"`
	Category        string            `json:"category"`
	Subcategory     string            `json:"subcategory"`
	DailyRate       float64           `json:"daily_rate"`
	WeeklyRate      float64           `json:"weekly_rate"`
	MonthlyRate     float64           `json:"monthly_rate"`
	SecurityDeposit float64           `json:"security_deposit"`
	Address         string            `json:"address"`
	City            string            `json:"city"`
	Latitude        float64           `json:"latitude"`
	Longitude       float64           `json:"longitude"`
	Specifications  map[string]string `json:"specifications"`
	Images          []string          `json:"images"`
}

func (h *HTTPHandler) HandleItems(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateItem(w, r)
	case http.MethodGet:
		itemID := r.URL.Query().Get("id")
		if itemID != "" {
			h.GetItem(w, r, itemID)
		} else {
			h.ListItems(w, r)
		}
	case http.MethodPut:
		h.UpdateItem(w, r)
	case http.MethodDelete:
		h.DeleteItem(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *HTTPHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var req CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ownerID, err := uuid.Parse(req.OwnerID)
	if err != nil {
		http.Error(w, "Invalid owner_id", http.StatusBadRequest)
		return
	}

	location := domain.Location{
		Address:   req.Address,
		City:      req.City,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	item, err := h.inventoryService.CreateItem(
		r.Context(), ownerID, req.Title, req.Description,
		domain.ItemCategory(req.Category), req.Subcategory,
		req.DailyRate, req.WeeklyRate, req.MonthlyRate, req.SecurityDeposit,
		location, req.Specifications, req.Images,
	)

	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         item.ID.String(),
		"title":      item.Title,
		"category":   item.Category,
		"daily_rate": item.DailyRate,
		"city":       item.City,
		"is_active":  item.IsActive,
		"created_at": item.CreatedAt,
	})
}

func (h *HTTPHandler) GetItem(w http.ResponseWriter, r *http.Request, itemID string) {
	id, err := uuid.Parse(itemID)
	if err != nil {
		http.Error(w, "Invalid item_id", http.StatusBadRequest)
		return
	}

	item, err := h.inventoryService.GetItem(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":               item.ID.String(),
		"owner_id":         item.OwnerID.String(),
		"title":            item.Title,
		"description":      item.Description,
		"category":         item.Category,
		"subcategory":      item.Subcategory,
		"daily_rate":       item.DailyRate,
		"weekly_rate":      item.WeeklyRate,
		"monthly_rate":     item.MonthlyRate,
		"security_deposit": item.SecurityDeposit,
		"city":             item.City,
		"address":          item.Address,
		"specifications":   item.Specifications,
		"images":           item.Images,
		"is_active":        item.IsActive,
		"created_at":       item.CreatedAt,
	})
}

func (h *HTTPHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 {
		pageSize = 20
	}

	category := r.URL.Query().Get("category")
	city := r.URL.Query().Get("city")
	minPriceStr := r.URL.Query().Get("min_price")
	maxPriceStr := r.URL.Query().Get("max_price")
	sort := r.URL.Query().Get("sort")

	filters := repository.ItemFilters{}
	if category != "" {
		cat := domain.ItemCategory(category)
		filters.Category = &cat
	}
	if city != "" {
		filters.City = &city
	}
	if minPriceStr != "" {
		if val, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filters.MinPrice = &val
		}
	}
	if maxPriceStr != "" {
		if val, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filters.MaxPrice = &val
		}
	}
	if sort != "" {
		filters.SortBy = &sort
	}

	items, total, err := h.inventoryService.ListItems(r.Context(), page, pageSize, filters)
	if err != nil {
		h.handleError(w, err)
		return
	}

	result := make([]map[string]interface{}, len(items))
	for i, item := range items {
		result[i] = map[string]interface{}{
			"id":         item.ID.String(),
			"title":      item.Title,
			"category":   item.Category,
			"city":       item.City,
			"daily_rate": item.DailyRate,
			"is_active":  item.IsActive,
			"images":     item.Images,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"items": result,
		"total": total,
		"page":  page,
	})
}

func (h *HTTPHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	itemID := r.URL.Query().Get("id")
	ownerID := r.URL.Query().Get("owner_id")

	id, err := uuid.Parse(itemID)
	if err != nil {
		http.Error(w, "Invalid item_id", http.StatusBadRequest)
		return
	}

	oid, err := uuid.Parse(ownerID)
	if err != nil {
		http.Error(w, "Invalid owner_id", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	item, err := h.inventoryService.UpdateItem(r.Context(), id, oid, updates)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":        item.ID.String(),
		"title":     item.Title,
		"is_active": item.IsActive,
	})
}

func (h *HTTPHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	itemID := r.URL.Query().Get("id")
	ownerID := r.URL.Query().Get("owner_id")

	id, err := uuid.Parse(itemID)
	if err != nil {
		http.Error(w, "Invalid item_id", http.StatusBadRequest)
		return
	}

	oid, err := uuid.Parse(ownerID)
	if err != nil {
		http.Error(w, "Invalid owner_id", http.StatusBadRequest)
		return
	}

	if err := h.inventoryService.DeleteItem(r.Context(), id, oid); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (h *HTTPHandler) GetOwnerItems(w http.ResponseWriter, r *http.Request) {
	ownerID := r.URL.Query().Get("owner_id")
	oid, err := uuid.Parse(ownerID)
	if err != nil {
		http.Error(w, "Invalid owner_id", http.StatusBadRequest)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	items, total, err := h.inventoryService.GetOwnerItems(r.Context(), oid, page, pageSize)
	if err != nil {
		h.handleError(w, err)
		return
	}

	result := make([]map[string]interface{}, len(items))
	for i, item := range items {
		result[i] = map[string]interface{}{
			"id":         item.ID.String(),
			"title":      item.Title,
			"category":   item.Category,
			"daily_rate": item.DailyRate,
			"is_active":  item.IsActive,
			"images":     item.Images,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"items": result,
		"total": total,
	})
}

func (h *HTTPHandler) SearchItems(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		query = r.URL.Query().Get("query")
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 {
		pageSize = 20
	}

	category := r.URL.Query().Get("category")
	city := r.URL.Query().Get("city")
	minPriceStr := r.URL.Query().Get("min_price")
	maxPriceStr := r.URL.Query().Get("max_price")
	sort := r.URL.Query().Get("sort")

	filters := repository.ItemFilters{}
	if category != "" {
		cat := domain.ItemCategory(category)
		filters.Category = &cat
	}
	if city != "" {
		filters.City = &city
	}
	if minPriceStr != "" {
		if val, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filters.MinPrice = &val
		}
	}
	if maxPriceStr != "" {
		if val, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filters.MaxPrice = &val
		}
	}
	if sort != "" {
		filters.SortBy = &sort
	}

	var items []*domain.RentalItem
	var total int
	var err error

	if query != "" {
		items, total, err = h.inventoryService.SearchItems(r.Context(), query, page, pageSize, filters)
	} else {
		items, total, err = h.inventoryService.ListItems(r.Context(), page, pageSize, filters)
	}

	if err != nil {
		h.handleError(w, err)
		return
	}

	result := make([]map[string]interface{}, len(items))
	for i, item := range items {
		result[i] = map[string]interface{}{
			"id":         item.ID.String(),
			"title":      item.Title,
			"category":   item.Category,
			"city":       item.City,
			"daily_rate": item.DailyRate,
			"is_active":  item.IsActive,
			"images":     item.Images,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"items": result,
		"total": total,
		"page":  page,
	})
}

func (h *HTTPHandler) BlockDates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ItemID    string `json:"item_id"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		BookingID string `json:"booking_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (h *HTTPHandler) CreateMaintenance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (h *HTTPHandler) handleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch err {
	case domain.ErrItemNotFound:
		w.WriteHeader(http.StatusNotFound)
	case domain.ErrUnauthorized:
		w.WriteHeader(http.StatusForbidden)
	case domain.ErrInvalidCategory:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
