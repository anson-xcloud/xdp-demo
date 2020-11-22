package server

import (
	"errors"
	"fmt"
)

var (
	// ErrAddressFormat address format invalid, must be appid:appsecret
	ErrAddressFormat = errors.New("address format invalid")

	// ErrTwiceWriteHTTPResponse twice write http response do not allowed
	ErrTwiceWriteHTTPResponse = errors.New("twice write http response")

	// ErrTimeout call timeout
	ErrTimeout = errors.New("timeout")

	ErrApiNowAllowed = errors.New("api donot allowed")

	ErrInvalidRemote = errors.New("invalid remote")

	ErrUnprepared = errors.New("server unprepared")
)

type StatusError struct {
	Status  int
	Message string
}

func NewStatusError(status int, message string) error {
	return &StatusError{
		Status:  status,
		Message: message,
	}
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("status code %d, %s", e.Status, e.Message)
}
