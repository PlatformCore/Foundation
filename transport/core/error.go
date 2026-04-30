package core

import "errors"

type Error struct {
	Code      string
	Message   string
	Status    int
	Temporary bool
	Cause     error
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	return e.Code
}
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}
func NewError(code, msg string, status int) *Error {
	return &Error{Code: code, Message: msg, Status: status}
}
func Wrap(code, msg string, status int, cause error) *Error {
	return &Error{Code: code, Message: msg, Status: status, Cause: cause}
}
func ToError(err error) *Error {
	if err == nil {
		return nil
	}
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return Wrap("INTERNAL", err.Error(), 500, err)
}
