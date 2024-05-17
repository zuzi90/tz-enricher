package models

import "errors"

var ErrUserNotFound = errors.New("not found")
var ErrNoRows = errors.New("err sql: no rows in result set")
