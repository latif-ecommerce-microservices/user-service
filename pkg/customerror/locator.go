package customerror

import (
	"fmt"
	"runtime"
	"strings"
)

type Locator struct {
	Line     int
	Function string
	File     string
}

func (el Locator) String() string {
	return fmt.Sprintf("%s:%d %s", el.File, el.Line, el.Function)
}

func WhereAmI(depthList ...int) *Locator {
	var depth int
	if depthList == nil {
		depth = 1
	} else {
		depth = depthList[0]
	}
	function, file, line, _ := runtime.Caller(depth)

	splitFile := strings.Split(file, "/")
	splitFunction := strings.Split(runtime.FuncForPC(function).Name(), "/")

	return &Locator{
		File:     splitFile[len(splitFile)-1],
		Function: splitFunction[len(splitFunction)-1],
		Line:     line,
	}
}
