package error

import "fmt"

type ErrNetwork string

// Error implements error.
func (e ErrNetwork) Error() string {
	return fmt.Sprintf("ErrNetwork: %s", string(e))
}

func NewErrNetwork(format string, args ...any) error {
	return ErrNetwork(fmt.Sprintf(format, args...))
}
