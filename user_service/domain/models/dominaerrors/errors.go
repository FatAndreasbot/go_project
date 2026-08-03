package dominaerrors

import "errors"

var UserNotFoundError error = errors.New("user not found")
var UserWrongPassword error = errors.New("wrong password")
