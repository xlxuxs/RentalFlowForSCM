package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rentalflow/auth-service/internal/domain"
	"github.com/rentalflow/auth-service/internal/repository"
	"github.com/rentalflow/auth-service/internal/token"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo    repository.UserRepository
	docRepo     repository.DocumentRepository
	jwtService  *token.JWTService
	passService *token.PasswordService
}

// NewAuthService creates a new auth service
func NewAuthService(
	userRepo repository.UserRepository,
	docRepo repository.DocumentRepository,
	jwtService *token.JWTService,
	passService *token.PasswordService,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		docRepo:     docRepo,
		jwtService:  jwtService,
		passService: passService,
	}
}

// AuthResult contains the result of authentication
type AuthResult struct {
	User         *domain.User
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, email, password, firstName, lastName, phone string, role domain.UserRole) (*AuthResult, error) {
	// Validate role
	if !role.IsValid() {
		return nil, domain.ErrInvalidRole
	}

	// Admin role can only be assigned by other admins
	if role == domain.RoleAdmin {
		return nil, domain.ErrForbidden
	}

	// Validate password strength
	if err := s.passService.ValidatePasswordStrength(password); err != nil {
		return nil, err
	}

	// Check if user already exists
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return nil, domain.ErrUserAlreadyExists
	}
	if err != domain.ErrUserNotFound {
		return nil, err
	}

	// Hash password
	passwordHash, err := s.passService.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := domain.NewUser(email, passwordHash, firstName, lastName, phone, role)

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate tokens
	tokenPair, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, err
	}

	// Store refresh token hash
	refreshHash := s.passService.HashRefreshToken(tokenPair.RefreshToken)
	if err := s.userRepo.UpdateRefreshToken(ctx, user.ID, refreshHash, &tokenPair.ExpiresAt); err != nil {
		return nil, err
	}

	return &AuthResult{
		User:         user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	// Verify password
	if !s.passService.VerifyPassword(password, user.PasswordHash) {
		return nil, domain.ErrInvalidCredentials
	}

	// Generate tokens
	tokenPair, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, err
	}

	// Store refresh token hash
	refreshHash := s.passService.HashRefreshToken(tokenPair.RefreshToken)
	if err := s.userRepo.UpdateRefreshToken(ctx, user.ID, refreshHash, &tokenPair.ExpiresAt); err != nil {
		return nil, err
	}

	return &AuthResult{
		User:         user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

// RefreshToken refreshes the access token using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*AuthResult, error) {
	// Hash the provided refresh token for future use
	_ = s.passService.HashRefreshToken(refreshToken)

	// We need to find the user with this refresh token
	// Since we're storing refresh token in user table, we need to iterate
	// In a real implementation, you might add an index or use Redis
	// For now, this is a simplified approach

	// This is a limitation of storing refresh token in user table
	// A better approach would be to decode user ID from the token itself
	// or use a separate refresh_tokens table

	return nil, domain.ErrInvalidRefreshToken
}

// RefreshTokenByUserID refreshes tokens for a known user
func (s *AuthService) RefreshTokenByUserID(ctx context.Context, userID uuid.UUID, refreshToken string) (*AuthResult, error) {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Verify refresh token
	providedHash := s.passService.HashRefreshToken(refreshToken)
	if user.RefreshTokenHash != providedHash {
		return nil, domain.ErrInvalidRefreshToken
	}

	// Check if refresh token is expired
	if user.RefreshTokenExpiresAt == nil || user.RefreshTokenExpiresAt.Before(time.Now()) {
		return nil, domain.ErrRefreshTokenExpired
	}

	// Generate new tokens
	tokenPair, err := s.jwtService.GenerateTokenPair(user)
	if err != nil {
		return nil, err
	}

	// Store new refresh token hash
	refreshHash := s.passService.HashRefreshToken(tokenPair.RefreshToken)
	if err := s.userRepo.UpdateRefreshToken(ctx, user.ID, refreshHash, &tokenPair.ExpiresAt); err != nil {
		return nil, err
	}

	return &AuthResult{
		User:         user,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}, nil
}

// Logout logs out a user
func (s *AuthService) Logout(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.ClearRefreshToken(ctx, userID)
}

// ValidateToken validates an access token
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*token.Claims, error) {
	return s.jwtService.ValidateAccessToken(tokenString)
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// UpdateProfile updates a user's profile
func (s *AuthService) UpdateProfile(ctx context.Context, userID uuid.UUID, firstName, lastName, phone, bio *string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if firstName != nil {
		user.FirstName = *firstName
	}
	if lastName != nil {
		user.LastName = *lastName
	}
	if phone != nil {
		user.Phone = *phone
	}
	if bio != nil {
		user.Bio = *bio
	}

	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateAvatar updates a user's avatar URL
func (s *AuthService) UpdateAvatar(ctx context.Context, userID uuid.UUID, avatarURL string) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.AvatarURL = avatarURL
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify currentpassword
	if !s.passService.VerifyPassword(currentPassword, user.PasswordHash) {
		return domain.ErrInvalidCredentials
	}

	// Validate new password strength
	if err := s.passService.ValidatePasswordStrength(newPassword); err != nil {
		return err
	}

	// Hash new password
	newPasswordHash, err := s.passService.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	user.PasswordHash = newPasswordHash
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(ctx, user)
}

// UploadDocument uploads an identity document
func (s *AuthService) UploadDocument(ctx context.Context, userID uuid.UUID, docType, docURL string) (*domain.IdentityDocument, error) {
	// Validate document type
	if docType != domain.DocTypeDriverLicense &&
		docType != domain.DocTypeNationalID &&
		docType != domain.DocTypePassport {
		return nil, domain.ErrInvalidDocumentType
	}

	// Verify user exists
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Create document
	doc := domain.NewIdentityDocument(userID, docType, docURL)
	if err := s.docRepo.Create(ctx, doc); err != nil {
		return nil, err
	}

	return doc, nil
}

// GetVerificationStatus gets a user's verification status
func (s *AuthService) GetVerificationStatus(ctx context.Context, userID uuid.UUID) (domain.VerificationStatus, []*domain.IdentityDocument, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", nil, err
	}

	docs, err := s.docRepo.GetByUserID(ctx, userID)
	if err != nil {
		return "", nil, err
	}

	return user.VerificationStatus, docs, nil
}

// VerifyUser updates a user's verification status (admin only)
func (s *AuthService) VerifyUser(ctx context.Context, userID uuid.UUID, status domain.VerificationStatus) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.VerificationStatus = status
	if status == domain.VerificationVerified {
		user.IdentityVerified = true
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// ListUsers lists users with pagination and filters (admin only)
func (s *AuthService) ListUsers(ctx context.Context, page, pageSize int, role *domain.UserRole, status *domain.VerificationStatus) ([]*domain.User, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	filters := repository.UserFilters{
		Role:               role,
		VerificationStatus: status,
	}

	return s.userRepo.List(ctx, offset, pageSize, filters)
}
