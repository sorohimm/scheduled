package errors

import "errors"

var ErrServerConfused = errors.New("server error")
var ErrBadRequest = errors.New("bad request")
var ErrInvalidGroup = errors.New("invalid group name")
var ErrGroupNotSet = errors.New("the group is not set")
var ErrNotFound = errors.New("not found")
