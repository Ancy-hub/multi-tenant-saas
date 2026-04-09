package utils

import "regexp"

func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}

func IsStrongPassword(password string)bool{
	if len(password)<8{
		return false
	}
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _,c:=range password{
		switch{
		case 'A'<=c && c<='Z':
			hasUpper=true
		case 'a'<=c && c<='z':
			hasLower=true
		case '0'<=c && c<='9':
			hasNumber=true
		default:
			hasSpecial=true
		}
	}
	return hasUpper && hasLower && hasNumber && hasSpecial
}