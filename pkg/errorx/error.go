package errorx

import (
	"errors"
)

type ErrorX struct {
	Code int
	error
}

func NewErrorX(code int, err any) *ErrorX {
	var e error
	switch err.(type) {
	case error:
		e = err.(error)
	case string:
		e = errors.New(err.(string))
	default:
		e = errors.New("unknown error")
	}
	return &ErrorX{
		code,
		e,
	}
}

func (e *ErrorX) GetCode() int {
	return e.Code
}

func (e *ErrorX) GetMsg() string {
	if e == nil || e.error == nil {
		return ""
	}
	return e.error.Error()
}

func (e *ErrorX) Error() string {
	return e.GetMsg()
}

func ParseErrorX(err error) *ErrorX {
	var errorX *ErrorX
	ok := errors.As(err, &errorX)
	if !ok {
		return NewErrorX(UnknownCode, errors.New("unknow error"))
	}
	return errorX
}
