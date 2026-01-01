package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/rentalflow/auth-service/internal/domain"
	"github.com/rentalflow/auth-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// AuthHandler implements the gRPC AuthService
type AuthHandler struct {
	UnimplementedAuthServiceServer
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register registers a new user
func (h *AuthHandler) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	// Validate required fields
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		return nil, status.Error(codes.InvalidArgument, "email, password, first_name, and last_name are required")
	}

	role := domain.UserRole(req.Role)
	if req.Role == "" {
		role = domain.RoleRenter
	}

	result, err := h.authService.Register(ctx, req.Email, req.Password, req.FirstName, req.LastName, req.Phone, role)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &AuthResponse{
		User:         toProtoUser(result.User),
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

// Login authenticates a user
func (h *AuthHandler) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	result, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &AuthResponse{
		User:         toProtoUser(result.User),
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}, nil
}

// RefreshToken refreshes the access token
func (h *AuthHandler) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*AuthResponse, error) {
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh_token is required")
	}

	// For simplicity, we'll implement this when we have the user context
	// In a real implementation, the refresh token would contain the user ID
	return nil, status.Error(codes.Unimplemented, "use RefreshTokenWithUserID")
}

// Logout logs out a user
func (h *AuthHandler) Logout(ctx context.Context, req *LogoutRequest) (*LogoutResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	if err := h.authService.Logout(ctx, userID); err != nil {
		return nil, toGRPCError(err)
	}

	return &LogoutResponse{Success: true}, nil
}

// ValidateToken validates an access token (internal service-to-service)
func (h *AuthHandler) ValidateToken(ctx context.Context, req *ValidateTokenRequest) (*ValidateTokenResponse, error) {
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	claims, err := h.authService.ValidateToken(ctx, req.Token)
	if err != nil {
		return &ValidateTokenResponse{Valid: false}, nil
	}

	return &ValidateTokenResponse{
		Valid:  true,
		UserId: claims.UserID,
		Role:   claims.Role,
		Email:  claims.Email,
	}, nil
}

// GetProfile retrieves the user's profile
func (h *AuthHandler) GetProfile(ctx context.Context, req *GetProfileRequest) (*User, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	user, err := h.authService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return toProtoUser(user), nil
}

// UpdateProfile updates the user's profile
func (h *AuthHandler) UpdateProfile(ctx context.Context, req *UpdateProfileRequest) (*User, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	var firstName, lastName, phone, bio *string
	if req.FirstName != nil {
		firstName = req.FirstName
	}
	if req.LastName != nil {
		lastName = req.LastName
	}
	if req.Phone != nil {
		phone = req.Phone
	}
	// Bio not in protobuf yet, pass nil
	bio = nil

	user, err := h.authService.UpdateProfile(ctx, userID, firstName, lastName, phone, bio)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return toProtoUser(user), nil
}

// GetUserById retrieves a user by ID (internal service-to-service)
func (h *AuthHandler) GetUserById(ctx context.Context, req *GetUserByIdRequest) (*User, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	user, err := h.authService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return toProtoUser(user), nil
}

// GetUserByIdAdmin retrieves a user by ID (admin only)
func (h *AuthHandler) GetUserByIdAdmin(ctx context.Context, req *GetUserByIdRequest) (*User, error) {
	return h.GetUserById(ctx, req)
}

// VerifyUser updates a user's verification status (admin only)
func (h *AuthHandler) VerifyUser(ctx context.Context, req *VerifyUserRequest) (*User, error) {
	if req.UserId == "" || req.Status == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and status are required")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	verificationStatus := domain.VerificationStatus(req.Status)
	user, err := h.authService.VerifyUser(ctx, userID, verificationStatus)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return toProtoUser(user), nil
}

// ListUsers lists users with pagination (admin only)
func (h *AuthHandler) ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error) {
	page := int(req.Page)
	pageSize := int(req.PageSize)

	var role *domain.UserRole
	if req.Role != nil && *req.Role != "" {
		r := domain.UserRole(*req.Role)
		role = &r
	}

	var status *domain.VerificationStatus
	if req.VerificationStatus != nil && *req.VerificationStatus != "" {
		s := domain.VerificationStatus(*req.VerificationStatus)
		status = &s
	}

	users, total, err := h.authService.ListUsers(ctx, page, pageSize, role, status)
	if err != nil {
		return nil, toGRPCError(err)
	}

	protoUsers := make([]*User, len(users))
	for i, user := range users {
		protoUsers[i] = toProtoUser(user)
	}

	return &ListUsersResponse{
		Users:    protoUsers,
		Total:    int32(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}

// UploadIdentityDocument uploads an identity document
func (h *AuthHandler) UploadIdentityDocument(ctx context.Context, req *UploadDocumentRequest) (*UploadDocumentResponse, error) {
	if req.UserId == "" || req.DocumentType == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and document_type are required")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	// In a real implementation, we would:
	// 1. Upload the document bytes to S3/MinIO
	// 2. Get the URL back
	// For now, we'll assume the document is already uploaded
	docURL := "https://storage.example.com/docs/" + uuid.New().String()

	doc, err := h.authService.UploadDocument(ctx, userID, req.DocumentType, docURL)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &UploadDocumentResponse{
		DocumentId:  doc.ID.String(),
		DocumentUrl: doc.DocumentURL,
	}, nil
}

// GetVerificationStatus gets the user's verification status
func (h *AuthHandler) GetVerificationStatus(ctx context.Context, req *GetVerificationStatusRequest) (*VerificationStatusResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	verificationStatus, docs, err := h.authService.GetVerificationStatus(ctx, userID)
	if err != nil {
		return nil, toGRPCError(err)
	}

	protoDocs := make([]*IdentityDocument, len(docs))
	for i, doc := range docs {
		protoDocs[i] = &IdentityDocument{
			Id:           doc.ID.String(),
			DocumentType: doc.DocumentType,
			DocumentUrl:  doc.DocumentURL,
			UploadedAt:   timestamppb.New(doc.UploadedAt),
		}
	}

	return &VerificationStatusResponse{
		Status:    string(verificationStatus),
		Documents: protoDocs,
	}, nil
}

// Helper functions

func toProtoUser(user *domain.User) *User {
	return &User{
		Id:                 user.ID.String(),
		Email:              user.Email,
		FirstName:          user.FirstName,
		LastName:           user.LastName,
		Phone:              user.Phone,
		Role:               string(user.Role),
		IdentityVerified:   user.IdentityVerified,
		VerificationStatus: string(user.VerificationStatus),
		CreatedAt:          timestamppb.New(user.CreatedAt),
		UpdatedAt:          timestamppb.New(user.UpdatedAt),
	}
}

func toGRPCError(err error) error {
	switch err {
	case domain.ErrUserNotFound:
		return status.Error(codes.NotFound, err.Error())
	case domain.ErrUserAlreadyExists:
		return status.Error(codes.AlreadyExists, err.Error())
	case domain.ErrInvalidCredentials:
		return status.Error(codes.Unauthenticated, err.Error())
	case domain.ErrInvalidToken, domain.ErrExpiredToken:
		return status.Error(codes.Unauthenticated, err.Error())
	case domain.ErrInvalidRefreshToken, domain.ErrRefreshTokenExpired:
		return status.Error(codes.Unauthenticated, err.Error())
	case domain.ErrUnauthorized:
		return status.Error(codes.Unauthenticated, err.Error())
	case domain.ErrForbidden:
		return status.Error(codes.PermissionDenied, err.Error())
	case domain.ErrInvalidRole, domain.ErrInvalidDocumentType:
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
