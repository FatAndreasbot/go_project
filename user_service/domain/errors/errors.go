package dominaerrors

import "errors"

var WrongPasswdOrNoUserFound error = errors.New("wrong password or user not found")
