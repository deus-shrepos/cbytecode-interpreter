package main

import (
	"errors"
	"fmt"
)

var ErrMalformedRequest = errors.New("Malformed Request")
var ErrInvalidPath = errors.New("Invalid HTTP Path")

type HttpError struct {
	Msg string
	Err error
}

func NewHttpError(err error, msg string) HttpError {
	return HttpError{
		Err: err,
		Msg: msg,
	}
}

func (e HttpError) Error() string {
	return fmt.Sprintf("%v: %v", e.Err, e.Msg)
}

func (e HttpError) Unwrap() error {
	return e.Err
}
