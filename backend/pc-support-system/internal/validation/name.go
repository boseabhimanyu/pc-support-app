package validation

import (
	"errors"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var nameRegex = regexp.MustCompile(`^[A-Za-zÀ-ÿ]+([ '-][A-Za-zÀ-ÿ]+)*$`)

var titleCaser = cases.Title(language.English)

func ValidateName(name string) (string, error) {

	name = strings.TrimSpace(name)
	name = titleCaser.String(strings.ToLower(name))

	if len(name) < 2 || len(name) > 50 {
		return "", errors.New("name must be between 2 and 50 characters")
	}

	if !nameRegex.MatchString(name) {
		return "", errors.New("name contains invalid characters")
	}

	return name, nil
}
