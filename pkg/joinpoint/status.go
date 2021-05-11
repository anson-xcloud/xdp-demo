package joinpoint

import (
	"errors"
	"fmt"
)

var (
	ErrUnreachable    = errors.New("")
	ErrProviderNeed   = errors.New("provider need")
	ErrDone           = errors.New("done")
	ErrNoneServerAddr = errors.New("no server addr")
)

var (
	CodeOK         = 0
	CodeUnkown     = 1
	CodeMaxReserve = 100
)

var (
	StatusOK = NewStatus(CodeOK, "ok")
)

type Status struct {
	Code    int
	Message string
	Details []interface{}
}

func FromError(err error) *Status {
	if err == nil {
		return StatusOK
	}

	if st, ok := err.(*Status); ok {
		return st
	}

	if s, ok := err.(interface {
		GetCode() int
		GetMessage() string
	}); ok {
		return NewStatus(s.GetCode(), s.GetMessage())
	}
	return NewStatus(CodeUnkown, err.Error())
}

// New returns a Status representing c and msg.
func NewStatus(c int, msg string) *Status {
	return &Status{Code: c, Message: msg}
}

// Newf returns New(c, fmt.Sprintf(format, a...)).
func NewStatusf(c int, format string, a ...interface{}) *Status {
	return NewStatus(c, fmt.Sprintf(format, a...))
}

func (s *Status) GetCode() int {
	return s.Code
}

func (s *Status) GetMessage() string {
	return s.Message
}

func (s *Status) Error() string {
	return fmt.Sprintf("status error code: %v, message: %s", s.Code, s.Message)
}

func (s *Status) WithDetail(details ...interface{}) *Status {
	s.Details = append(s.Details, details...)
	return s
}
