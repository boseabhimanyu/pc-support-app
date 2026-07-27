package validation

import (
	"errors"
	"regexp"
	"strings"
)

var nameRegex = regexp.MustCompile(`^[A-Za-zÀ-ÿ]+([ '-][A-Za-zÀ-ÿ]+)*$`)

func ValidateName(name string) (string, error) {

	name = strings.TrimSpace(name)

	if len(name) < 2 || len(name) > 50 {
		return "", errors.New("name must be between 2 and 50 characters")
	}

	if !nameRegex.MatchString(name) {
		return "", errors.New("name contains invalid characters")
	}

	return name, nil
}
