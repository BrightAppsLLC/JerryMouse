package Utils

import (
	"path"
	"runtime"
)

// GetPathForThisPackage -
func GetPathForThisPackage() string {
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		return path.Dir(filename) + "/"
	}

	return "./"
}
