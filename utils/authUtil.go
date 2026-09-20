package utils

import "regexp"

func IsPasswordValid(password string) bool {
	var (
		hasMinLen     = len(password) >= 8
		hasUpper, _   = regexp.MatchString(`[A-Z]`, password)
		hasLower, _   = regexp.MatchString(`[a-z]`, password)
		hasSpecial, _ = regexp.MatchString(`[\W_]`, password)
	)
	return hasMinLen && hasUpper && hasLower && hasSpecial
}