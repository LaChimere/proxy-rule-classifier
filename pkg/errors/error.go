package errors

import (
	"fmt"
	"log"
)

type Error struct {
	code int
	msg  string
}

var codes = map[int]string{}

func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		log.Panicf("error code %d already exists", code)
	}

	codes[code] = msg
	return &Error{code: code, msg: msg}
}

func (err *Error) Error() string {
	return fmt.Sprintf("error code: %d, error message: %s", err.code, err.msg)
}

func (err *Error) Code() int {
	return err.code
}

func (err *Error) Msg() string {
	return err.msg
}
