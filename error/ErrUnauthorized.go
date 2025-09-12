package error

import "fmt"

type ErrUnauthorized string

// Error implements error.
func (e ErrUnauthorized) Error() string {
	return fmt.Sprintf("ErrUnauthorized: %s", string(e))
}

func NewErrUnauthorized(format string, args ...any) error {
	return ErrUnauthorized(fmt.Sprintf(format, args...))
}
