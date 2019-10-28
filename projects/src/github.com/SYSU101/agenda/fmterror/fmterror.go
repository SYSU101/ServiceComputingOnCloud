package fmterror

import (
	"errors"
	"fmt"
)

func New(format string, v ...interface{}) error {
	return errors.New(fmt.Sprintf(format, v...))
}
