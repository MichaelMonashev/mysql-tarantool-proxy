// Используется для вывода отлаточных сообщений и сообщений об ошибках.
package warn

import (
	"errors"
	"fmt"
	"os"
)

func Errorln(a ...interface{}) error {
	return errors.New(fmt.Sprintln(a...))
}

func Warn(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}
