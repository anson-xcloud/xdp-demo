package xdp

import "errors"

var (
	// ErrAddressFormat address format invalid, must be appid:appsecret
	ErrAddressFormat = errors.New("address format invalid")

	// ErrTwiceWriteHTTPResponse twice write http response do not allowed
	ErrTwiceWriteHTTPResponse = errors.New("twice write http response")

	// ErrTimeout call timeout
	ErrTimeout = errors.New("timeout")
)
