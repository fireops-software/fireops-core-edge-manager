package error

import "fmt"

type ErrChannelClosed string

// Error implements error.
func (e ErrChannelClosed) Error() string {
	return fmt.Sprintf("ErrChannelClosed: %s", string(e))
}

func NewErrChannelClosed(format string, args ...any) error {
	return ErrChannelClosed(fmt.Sprintf(format, args...))
}
