package formaterror

import (
	"errors"
	"strings"
)

func FormatError(err string) error {
	if strings.Contains(err, "nickname") {
		return errors.New("Nickname already taken")
	}
	if strings.Contains(err, "email") {
		return errors.New("Email already taken")
	}

	if strings.Contains(err, "title") {
		return errors.New("Title already taken")
	}

	if strings.Contains(err, "hashedPassword") {
		return errors.New("Incorrect password")
	}

	return errors.New("Incorrect details")
}