package utils

import (
	"strconv"
	"strings"
)

const (
	iranNationalCodeLen = 10
)

func ValidateNationalCode(code string) bool {
	if code != "" {
		var err error
		if len([]rune(code)) != iranNationalCodeLen {
			return false //National Code must be exactly 10 digits
		}

		digits := strings.Split(code, "")
		sum := 0

		for i := 0; i < 9; i++ {
			digit, err := strconv.Atoi(digits[i])
			if err != nil {
				return false
			}

			sum += digit * (10 - i)
		}

		var chkSum int
		r := sum % 11

		if r < 2 {
			chkSum = r
		} else {
			chkSum = 11 - (r)
		}

		ld, err := strconv.Atoi(digits[9])
		if err != nil {
			return false
		}

		if ld != chkSum {
			return false // Invalid National Code
		}

		return true
	}

	return false
}
