package formaterror

import (
	"errors"
	"log"
	"strings"
)

func FormatError(err string) error {
	log.Println(err)
	if strings.Contains(err, "nickname") {
		return errors.New("Nickname already taken")
	}
	if strings.Contains(err, "email") {
		return errors.New("Email already taken")
	}

	if strings.Contains(err, "title") {
		return errors.New("title already taken")
	}

	if strings.Contains(err, "hashedPassword") {
		return errors.New("Incorrect password")
	}

	return errors.New("Incorrect details")
}