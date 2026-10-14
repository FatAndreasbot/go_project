package dominaerrors

import "errors"

var (
	InvalidCredentials = errors.New("invalid credentials")
	ExpiredCredentials = errors.Join(InvalidCredentials, errors.New("expired credentials"))

	NotEnoughPermission = errors.New("lacking permissions")
)
