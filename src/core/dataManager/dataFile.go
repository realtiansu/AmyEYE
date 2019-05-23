package dataManager

import (
	"core/common/errors"
	"fmt"
	"os"
)

func record(line string, path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return errors.New("can't open file %s ", path)
	}
	_, _ = fmt.Fprintln(file, line)

	return err
}
