package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rentalflow/review-service/internal/domain"
	"github.com/rentalflow/review-service/internal/repository"
)

type ReviewService struct {
	reviewRepo repository.ReviewRepository
}

func NewReviewService(reviewRepo repository.ReviewRepository) *ReviewService {
	return &ReviewService{reviewRepo: reviewRepo}
}

func (s *ReviewService) CreateReview(ctx context.Context, itemID uuid.UUID, bookingID *uuid.UUID, reviewerID uuid.UUID, reviewType domain.ReviewType, rating float64, comment string) (*domain.Review, error) {
	if rating < 1.0 || rating > 5.0 {
		return nil, domain.ErrInvalidRating
	}

	// Use bookingID if provided, otherwise use nil UUID
	var bid uuid.UUID
	if bookingID != nil {
		bid = *bookingID
	} else {
		bid = uuid.Nil
	}

	review := domain.NewReview(bid, reviewerID, reviewType, rating, comment)

	// Set target IDs based on type
	if reviewType == domain.TypeRenterToItem {
		review.TargetItemID = &itemID
	} else if reviewType == domain.TypeRenterToOwner || reviewType == domain.TypeOwnerToRenter {
		review.TargetUserID = &itemID // In this service, itemID is used as target ID for users too in the handler's call
	}

	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return nil, err
	}
	return review, nil
}

func (s *ReviewService) GetReview(ctx context.Context, reviewID uuid.UUID) (*domain.Review, error) {
	return s.reviewRepo.GetByID(ctx, reviewID)
}

func (s *ReviewService) GetItemReviews(ctx context.Context, itemID uuid.UUID, page, pageSize int) ([]*domain.Review, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.reviewRepo.GetByItem(ctx, itemID, offset, pageSize)
}

func (s *ReviewService) GetUserReviews(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*domain.Review, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.reviewRepo.GetByUser(ctx, userID, offset, pageSize)
}

func (s *ReviewService) UpdateReview(ctx context.Context, reviewID uuid.UUID, rating float64, comment string) (*domain.Review, error) {
	review, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return nil, err
	}

	if rating >= 1.0 && rating <= 5.0 {
		review.Rating = rating
	}
	if comment != "" {
		review.Comment = comment
	}

	if err := s.reviewRepo.Update(ctx, review); err != nil {
		return nil, err
	}
	return review, nil
}

func (s *ReviewService) DeleteReview(ctx context.Context, reviewID uuid.UUID) error {
	return s.reviewRepo.Delete(ctx, reviewID)
}
