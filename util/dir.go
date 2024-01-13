package util

import (
	"path"
	"runtime"
)

func GetRootPath() string {
	_, filename, _, _ := runtime.Caller(0)
	RootPath := path.Dir(path.Dir(filename))
	return RootPath
}
