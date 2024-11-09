//go:build !solution

package spacecollapse

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func CollapseSpaces(input string) string {
	var buf strings.Builder
	inSpace := false

	for _, c := range input {
		switch {
		case !utf8.ValidRune(c):
			buf.WriteRune('\uFFFD')
		case unicode.IsSpace(c):
			if !inSpace {
				buf.WriteRune(' ')
			}
			inSpace = true
		default:
			buf.WriteRune(c)
			inSpace = false
		}
	}

	return buf.String()
}
