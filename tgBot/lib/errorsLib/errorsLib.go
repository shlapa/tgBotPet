package errorsLib

import (
	"errors"
	"fmt"
)

var ErrNoSavedPage = errors.New("no saved page")
var ErrUnknownProcess = errors.New("unknown process")
var ErrorTypeMeta = errors.New("error type meta")

func Wrap(m string, err error) error {
	if err == nil {
		return nil
	} else {
		return fmt.Errorf("%s %w", m, err)
	}
}
