package rule

type Error struct {
	msg string
}

func NewError(msg string) *Error {
	return &Error{msg: msg}
}

func (err *Error) Error() string {
	return err.msg
}
