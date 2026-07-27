package validation

import (
	"errors"
	"regexp"
	"strings"
)

var phoneRegex = regexp.MustCompile(`^[0-9]{10}$`)

func ValidatePhone(phone string) (string, error) {

	phone = strings.TrimSpace(phone)

	if !phoneRegex.MatchString(phone) {
		return "", errors.New("phone number must contain exactly 10 digits")
	}

	return phone, nil
}
