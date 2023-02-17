package validators

import "regexp"

func IsValidPassword(password string) bool {
	if len(password) < 8 || len(password) > 20 {
		return false
	}

	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false
	for _, char := range password {
		switch {
		case 'a' <= char && char <= 'z':
			hasLower = true
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case '0' <= char && char <= '9':
			hasDigit = true
		}
	}

	specialCharRegex := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`)
	hasSpecial = specialCharRegex.MatchString(password)

	return hasLower && hasUpper && hasDigit && hasSpecial
}
