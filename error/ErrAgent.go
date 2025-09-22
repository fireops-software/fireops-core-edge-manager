package error

import "fmt"

type ErrAgent string

// Error implements error.
func (e ErrAgent) Error() string {
	return fmt.Sprintf("ErrAgent: %s", string(e))
}

func NewErrAgent(format string, args ...any) error {
	return ErrAgent(fmt.Sprintf(format, args...))
}
