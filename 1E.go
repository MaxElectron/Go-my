//go:build !solution

package speller

import (
	"fmt"
)

func Spell(n int64) string {
	ans := ""

	if n == 0 {
		return "zero"
	}

	// if negative add "minus "
	if n < 0 {
		ans += "minus "
		n = -n
	}

	// convert n to string
	num := fmt.Sprintf("%d", n)

	// handle billions if any
	if len(num) > 9 {
		ans += Spell(n/1000000000) + " billion"
		num = fmt.Sprintf("%d", n%1000000000)
		if n%1000000000 != 0 {
			ans += " "
		}
	}

	// handle millions if any
	if len(num) > 6 {
		ans += Spell(n%1000000000/1000000) + " million"
		num = fmt.Sprintf("%d", n%1000000)
		if n%1000000 != 0 {
			ans += " "
		}
	}

	// handle thousands if any
	if len(num) > 3 {
		ans += Spell(n%1000000/1000) + " thousand"
		num = fmt.Sprintf("%d", n%1000)
		if n%1000 != 0 {
			ans += " "
		}
	}

	// handle hundreds if any
	if len(num) > 2 {
		ans += Spell(n%1000/100) + " hundred"
		num = fmt.Sprintf("%d", n%100)
		if n%100 != 0 {
			ans += " "
		}
	}

	// handle the case when the remainder is at least 20
	if n%100 >= 20 {
		switch num[0] {
		case '2':
			ans += "twenty"
		case '3':
			ans += "thirty"
		case '4':
			ans += "forty"
		case '5':
			ans += "fifty"
		case '6':
			ans += "sixty"
		case '7':
			ans += "seventy"
		case '8':
			ans += "eighty"
		case '9':
			ans += "ninety"
		}
		if num[1] != '0' {
			ans += "-"
		}
		num = fmt.Sprintf("%d", n%10)
	}

	// handle one through nineteen
	switch num {
	case "1":
		ans += "one"
	case "2":
		ans += "two"
	case "3":
		ans += "three"
	case "4":
		ans += "four"
	case "5":
		ans += "five"
	case "6":
		ans += "six"
	case "7":
		ans += "seven"
	case "8":
		ans += "eight"
	case "9":
		ans += "nine"
	case "10":
		ans += "ten"
	case "11":
		ans += "eleven"
	case "12":
		ans += "twelve"
	case "13":
		ans += "thirteen"
	case "14":
		ans += "fourteen"
	case "15":
		ans += "fifteen"
	case "16":
		ans += "sixteen"
	case "17":
		ans += "seventeen"
	case "18":
		ans += "eighteen"
	case "19":
		ans += "nineteen"
	}

	return ans
}
