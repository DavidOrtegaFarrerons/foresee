package models

import "errors"

var (
	ErrInvalidCredentials           = errors.New("models: invalid credentials")
	ErrEmailAlreadyExists           = errors.New("models: email already exists")
	ErrUsernameAlreadyExists        = errors.New("models: username already exists")
	ErrUserAlreadyBetOnMarket       = errors.New("user has already placed a bet on this market")
	ErrUserNotAuthorized            = errors.New("user not authorized to do the following operation")
	ErrMarketAlreadyResolved        = errors.New("this market has already been resolved")
	ErrMarketNotExpired             = errors.New("market has not expired yet")
	ErrOutcomeDoesNotBelongToMarket = errors.New("this outcome does not belong to this market")
)
