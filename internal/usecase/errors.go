package usecase

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrNotImplemented      = errors.New("not implemented")
	ErrPaymentRequired     = errors.New("payment required")
	ErrConflict            = errors.New("conflict")
	ErrUnprocessableEntity = errors.New("unprocessable entity")
)
