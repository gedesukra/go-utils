package goutils

import (
	"fmt"
	"runtime"
)

var ConsUnableReadBody string = "Unable read http body"
var ConsUnableUnmarshal string = "Unable unmarshal json"
var ConsUnableMarshal string = "Unable marshal obj"

func getStacktraceError(err error) string {
	pc, fn, line, _ := runtime.Caller(1)
	return fmt.Sprintf("stacktrace -> in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
}
