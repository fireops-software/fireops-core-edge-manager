package error

import "fmt"

type ErrNotImplemented string

// Error implements error.
func (e ErrNotImplemented) Error() string {
	return fmt.Sprintf("ErrNotImplemented: %s", string(e))
}

func NewErrNotImplemented(format string, args ...any) error {
	return ErrNotImplemented(fmt.Sprintf(format, args...))
}
