package Utils

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"strings"
)

var mainPackagePath string

// GetPathForThisPackage -
func GetPathForThisPackage() string {
	_, filename, _, ok := runtime.Caller(1)
	if ok {
		return path.Dir(filename) + "/"
	}

	return "./"
}

// InitMainPackagePath - to be called in `main()`
func InitMainPackagePath() {
	mainPackagePath = GetPathForThisPackage()
}

// CallStack - returns function name at runtime
func CallStack() string {
	pc, _, _, _ := runtime.Caller(1)
	var functionName = runtime.FuncForPC(pc).Name()

	return strings.Replace(functionName, "_"+mainPackagePath, "", 1)
}

// PackageName - returns function name at runtime
func PackageName(i interface{}) string {
	var packagePath = reflect.TypeOf(i).PkgPath()
	runes := []rune(packagePath)
	return string(runes[1:len(packagePath)])
}

// Trace - traces a call with line number
func Trace() {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	fmt.Printf("%s,:%d %s\n", frame.File, frame.Line, frame.Function)
}
