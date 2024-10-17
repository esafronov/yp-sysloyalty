package helpers

import (
	"regexp"
)

func ValidateOnlyDigits(str string) bool {

	var re = regexp.MustCompile(`^[0-9]+$`)

	if re.MatchString(str) {
		return true
	} else {
		return false
	}

}
