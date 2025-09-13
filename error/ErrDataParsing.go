package error

import "fmt"

type ErrDataParsing string

// Error implements error.
func (e ErrDataParsing) Error() string {
	return fmt.Sprintf("ErrDataParsing: %s", string(e))
}

func NewErrDataParsing(format string, args ...any) error {
	return ErrDataParsing(fmt.Sprintf(format, args...))
}
