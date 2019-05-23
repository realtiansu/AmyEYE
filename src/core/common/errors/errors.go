package errors

import (
	"errors"
	"fmt"
)

func New(err string, value ...interface{}) error {
	msg := fmt.Sprintf(err, value...)
	fmt.Println(msg)
	return errors.New(msg)
}
