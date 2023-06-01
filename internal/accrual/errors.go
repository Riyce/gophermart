package accrual

import "errors"

var errNoContent error = errors.New("no content")
var errTooManyRequests error = errors.New("too many requests")
var errServerError error = errors.New("server error")
var errWrongStatusCode error = errors.New("wrong status code")
