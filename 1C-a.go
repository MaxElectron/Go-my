//go:build !solution

package reverse

func Reverse(input string) string {
	runes := []rune(input)
	left, right := 0, len(runes)-1
	for left < right {
		runes[left], runes[right] = runes[right], runes[left]
		left, right = left+1, right-1
	}
	return string(runes)
}

