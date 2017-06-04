// +build !release

package warn

import (
	"fmt"
	"os"
	"runtime/debug"
)

func Fatal(a ...interface{}) {
	Warn(a...)
	debug.PrintStack()
	os.Exit(1)
}

// Выводит отладочное сообщение в STDOUT
func Say(a ...interface{}) {
	fmt.Println(a...)
}
func SayS(a ...interface{}) {
	Say(a...)
	debug.PrintStack()
}
