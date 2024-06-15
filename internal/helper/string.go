package helper

import "strings"

func IsEmpty(s string) bool {
	return strings.Trim(s, " ") == ""
}
