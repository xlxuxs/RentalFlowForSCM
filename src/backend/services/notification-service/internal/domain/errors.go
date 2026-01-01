package domain

import "errors"

var (
	ErrNotificationNotFound = errors.New("notification not found")
	ErrMessageNotFound      = errors.New("message not found")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrInvalidChannel       = errors.New("invalid notification channel")
)
