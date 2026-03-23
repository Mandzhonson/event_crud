package apperr

import "errors"

var (
	BadRequest       = errors.New("bad request")
	InvalidReqParams = errors.New("invalid request parameters")
	EventNotFound    = errors.New("event not found")
	InternalServErr  = errors.New("internal server error")
	ErrTimeout       = errors.New("timeout")
	ErrCancel        = errors.New("context cancelled")
)
