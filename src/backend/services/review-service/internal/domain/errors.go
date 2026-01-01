package domain

import "errors"

var (
	ErrReviewNotFound = errors.New("review not found")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrInvalidRating  = errors.New("invalid rating value")
)
