package domainerrors

import "errors"

var (
	InvalidCredentials = errors.New("invalid credentials")
	ExpiredCredentials = errors.Join(InvalidCredentials, errors.New("expired credentials"))

	NotEnoughPermission = errors.New("lacking permissions")

	NotImplemented = errors.New("method was not implemented")
	NotFound       = errors.New("resource not found")
)
