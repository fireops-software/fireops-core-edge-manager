package error

import "fmt"

type ErrConflict string

// Error implements error.
func (e ErrConflict) Error() string {
	return fmt.Sprintf("ErrConflict: %s", string(e))
}

func NewErrConflict(format string, args ...any) error {
	return ErrConflict(fmt.Sprintf(format, args...))
}
