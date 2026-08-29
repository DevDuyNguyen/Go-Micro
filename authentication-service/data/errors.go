package data

import "errors"

var ErrPasswordNotMatch = errors.New("Password not match")
var ErrEntityNotFound = errors.New("Entity not found")
var ErrCanNotConnectToDB= errors.New("Can't connect to the database")