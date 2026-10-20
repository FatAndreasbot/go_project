package sessionstorage

import "errors"

var (
	DoesNotExists = errors.New("session does not exists")
	Expired       = errors.New("session was expired")
)
