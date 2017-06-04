// +build release

package warn

import (
	"fmt"
	"os"
)

func Fatal(a ...interface{}) {
	Warn(a...)
	os.Exit(1)
}

// печатает сообщение, которое выводится только при отладке, и потому в релизе отсуствующее
func Say(a ...interface{}) {
}

func SayS(a ...interface{}) {
}
