//go:build !solution

package varfmt

import (
	"bytes"
	"fmt"
	"strconv"
)

func Sprintf(format string, args ...interface{}) string {
	// create buffer and args counter
	currentArg := 0
	buffer := bytes.Buffer{}

	// create map for storing converted args, for optimization
	argValues := make(map[int]string)

	// iterate over format string
	for idx := 0; idx < len(format); idx++ {
		// if current char is not a placeholder
		if format[idx] != '{' {
			// just add it to buffer
			buffer.WriteByte(format[idx])
		} else {
			// infer the arg index
			argIndex := currentArg
			// update idx
			idx++
			// read the arg index if not empty
			if format[idx] != '}' {
				initialIdx := idx
				for format[idx] != '}' {
					idx++
				}
				argIndex, _ = strconv.Atoi(format[initialIdx:idx])
			}

			// if corresponding arg has not been converted yet convert it
			if _, exists := argValues[argIndex]; !exists {
				argValues[argIndex] = fmt.Sprint(args[argIndex])
			}

			// add converted arg to buffer
			buffer.WriteString(argValues[argIndex])
			// update current arg
			currentArg++
		}
	}

	return buffer.String()
}
