package error

import "fmt"

type ErrNotFound string

// Error implements error.
func (e ErrNotFound) Error() string {
	return fmt.Sprintf("ErrNotFound: %s", string(e))
}

func NewErrNotFound(format string, args ...any) error {
	return ErrNotFound(fmt.Sprintf(format, args...))
}
