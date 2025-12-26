package utils

import "regexp"

func IsPasswordStrong(pass string) bool {
	if len(pass) < 8 {
		return false
	}
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(pass)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(pass)
	hasSpecial := regexp.MustCompile(`[!@#~$%^&*()+|_]`).MatchString(pass)
	return hasNumber && hasUpper && hasSpecial
}
