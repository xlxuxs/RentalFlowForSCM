package domain

import "errors"

// Domain errors
var (
	ErrBookingNotFound     = errors.New("booking not found")
	ErrUnauthorized        = errors.New("unauthorized to perform this action")
	ErrInvalidStatus       = errors.New("invalid booking status")
	ErrInvalidDates        = errors.New("invalid booking dates")
	ErrDateConflict        = errors.New("booking dates conflict with existing reservation")
	ErrAlreadyCancelled    = errors.New("booking is already cancelled")
	ErrCannotCancel        = errors.New("booking cannot be cancelled")
	ErrAgreementNotSigned  = errors.New("rental agreement not signed")
	ErrPaymentNotCompleted = errors.New("payment not completed")
)
