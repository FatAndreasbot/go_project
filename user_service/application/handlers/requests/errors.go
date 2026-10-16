package requests

import "errors"

var (
	DBError      = errors.New("DB query error")
	JWTError     = errors.New("token encoding/decoding error")
	SessionError = errors.New("session error")
)
