package formaterror

import (
	"errors"
	"strings"
)

func FormatError(err string) error {
	if strings.Contains(err, "nickname") {
		return errors.New("[ERROR] Nickname already taken")
	}

	if strings.Contains(err, "email") {
		return errors.New("[ERROR] Email already taken")
	}

	if strings.Contains(err, "title") {
		return errors.New("[ERROR] Title already taken")
	}

	if strings.Contains(err, "hashedPassword") {
		return errors.New("[ERROR] Incorrect password")
	}

	return errors.New("Incorrect details")
}