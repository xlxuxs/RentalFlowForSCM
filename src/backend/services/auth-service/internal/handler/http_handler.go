package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/rentalflow/auth-service/internal/domain"
	"github.com/rentalflow/auth-service/internal/service"
)

// HTTPHandler provides REST endpoints for testing
type HTTPHandler struct {
	authService *service.AuthService
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(authService *service.AuthService) *HTTPHandler {
	return &HTTPHandler{authService: authService}
}

// RegisterRoutes registers HTTP routes
func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/ready", h.Ready)
	mux.HandleFunc("/api/auth/register", h.Register)
	mux.HandleFunc("/api/auth/login", h.Login)
	mux.HandleFunc("/api/auth/logout", h.Logout)
	mux.HandleFunc("/api/auth/profile", h.ProfileHandler)
	mux.HandleFunc("/api/auth/avatar", h.UpdateAvatar)
	mux.HandleFunc("/api/auth/change-password", h.ChangePassword)
	mux.HandleFunc("/api/auth/validate", h.ValidateToken)
	mux.HandleFunc("/api/users", h.ListUsers)
}

// Health check
func (h *HTTPHandler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

// Ready check
func (h *HTTPHandler) Ready(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}

// RegisterRequest for HTTP API
type RegisterHTTPRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
}

// Register handles user registration
func (h *HTTPHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterHTTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	role := domain.UserRole(req.Role)
	if req.Role == "" {
		role = domain.RoleRenter
	}

	result, err := h.authService.Register(r.Context(), req.Email, req.Password, req.FirstName, req.LastName, req.Phone, role)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": map[string]interface{}{
			"id":         result.User.ID.String(),
			"email":      result.User.Email,
			"first_name": result.User.FirstName,
			"last_name":  result.User.LastName,
			"role":       result.User.Role,
		},
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"expires_in":    result.ExpiresIn,
	})
}

// LoginRequest for HTTP API
type LoginHTTPRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login handles user authentication
func (h *HTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginHTTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": map[string]interface{}{
			"id":         result.User.ID.String(),
			"email":      result.User.Email,
			"first_name": result.User.FirstName,
			"last_name":  result.User.LastName,
			"role":       result.User.Role,
		},
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"expires_in":    result.ExpiresIn,
	})
}

// Logout handles user logout
func (h *HTTPHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Invalid user_id format", http.StatusBadRequest)
		return
	}

	if err := h.authService.Logout(r.Context(), uid); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// ProfileHandler handles both GET and PUT for user profile
func (h *HTTPHandler) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.GetProfile(w, r)
	} else if r.Method == "PUT" {
		h.UpdateProfile(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetProfile retrieves user profile
func (h *HTTPHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		token := h.extractToken(r)
		if token != "" {
			claims, err := h.authService.ValidateToken(r.Context(), token)
			if err == nil {
				userID = claims.UserID
			}
		}
	}

	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "Invalid user_id format", http.StatusBadRequest)
		return
	}

	user, err := h.authService.GetUserByID(r.Context(), uid)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":                  user.ID.String(),
		"email":               user.Email,
		"first_name":          user.FirstName,
		"last_name":           user.LastName,
		"phone":               user.Phone,
		"bio":                 user.Bio,
		"avatar_url":          user.AvatarURL,
		"role":                user.Role,
		"identity_verified":   user.IdentityVerified,
		"verification_status": user.VerificationStatus,
		"created_at":          user.CreatedAt,
	})
}

// UpdateProfile updates user profile
func (h *HTTPHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID    string  `json:"user_id"`
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Phone     *string `json:"phone"`
		Bio       *string `json:"bio"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "Invalid user_id format", http.StatusBadRequest)
		return
	}

	user, err := h.authService.UpdateProfile(r.Context(), uid, req.FirstName, req.LastName, req.Phone, req.Bio)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         user.ID.String(),
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"phone":      user.Phone,
		"bio":        user.Bio,
		"avatar_url": user.AvatarURL,
		"role":       user.Role,
	})
}

// UpdateAvatar updates user avatar
func (h *HTTPHandler) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID    string `json:"user_id"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.AvatarURL == "" {
		http.Error(w, "user_id and avatar_url are required", http.StatusBadRequest)
		return
	}

	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "Invalid user_id format", http.StatusBadRequest)
		return
	}

	user, err := h.authService.UpdateAvatar(r.Context(), uid, req.AvatarURL)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         user.ID.String(),
		"avatar_url": user.AvatarURL,
		"success":    true,
	})
}

// ChangePassword changes user password
func (h *HTTPHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID          string `json:"user_id"`
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.CurrentPassword == "" || req.NewPassword == "" {
		http.Error(w, "user_id, current_password, and new_password are required", http.StatusBadRequest)
		return
	}

	uid, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "Invalid user_id format", http.StatusBadRequest)
		return
	}

	if err := h.authService.ChangePassword(r.Context(), uid, req.CurrentPassword, req.NewPassword); err != nil {
		h.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

// ValidateToken validates an access token
func (h *HTTPHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	claims, err := h.authService.ValidateToken(r.Context(), req.Token)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"valid": false,
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":   true,
		"user_id": claims.UserID,
		"email":   claims.Email,
		"role":    claims.Role,
	})
}

// ListUsers lists all users (admin)
func (h *HTTPHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, total, err := h.authService.ListUsers(r.Context(), 1, 20, nil, nil)
	if err != nil {
		h.handleError(w, err)
		return
	}

	userList := make([]map[string]interface{}, len(users))
	for i, u := range users {
		userList[i] = map[string]interface{}{
			"id":         u.ID.String(),
			"email":      u.Email,
			"first_name": u.FirstName,
			"last_name":  u.LastName,
			"role":       u.Role,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users": userList,
		"total": total,
	})
}

func (h *HTTPHandler) handleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch err {
	case domain.ErrUserNotFound:
		w.WriteHeader(http.StatusNotFound)
	case domain.ErrUserAlreadyExists:
		w.WriteHeader(http.StatusConflict)
	case domain.ErrInvalidCredentials:
		w.WriteHeader(http.StatusUnauthorized)
	case domain.ErrForbidden:
		w.WriteHeader(http.StatusForbidden)
	case domain.ErrInvalidRole, domain.ErrInvalidDocumentType:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
func (h *HTTPHandler) extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}
