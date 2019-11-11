package mock

import (
	"path"
	"runtime"
)

func Getwd() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	return path.Dir(filename)
}
