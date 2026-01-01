package domain

import "errors"

var (
	ErrPaymentNotFound      = errors.New("payment not found")
	ErrInvoiceNotFound      = errors.New("invoice not found")
	ErrInvalidAmount        = errors.New("invalid payment amount")
	ErrPaymentFailed        = errors.New("payment processing failed")
	ErrRefundNotAllowed     = errors.New("refund not allowed for this payment")
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
	ErrUnauthorized         = errors.New("unauthorized to perform this action")
)
