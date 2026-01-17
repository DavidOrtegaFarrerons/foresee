package models

import "errors"

var (
	ErrInvalidCredentials     = errors.New("models: invalid credentials")
	ErrEmailAlreadyExists     = errors.New("models: email already exists")
	ErrUsernameAlreadyExists  = errors.New("models: username already exists")
	ErrUserAlreadyBetOnMarket = errors.New("user has already placed a bet on this market")
)
