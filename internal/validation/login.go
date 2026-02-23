package validation

import (
	"errors"
	"strings"
)

func LoginValidation(email, password string) error {
	if err := notEmpty(email, password); err != nil {
		return err
	}
	if err := passwordLen(password); err != nil {
		return err
	}
	if err := emailValidate(email); err != nil {
		return err
	}

	return nil
}

func notEmpty(email, password string) error {

	if email == "" {
		return errors.New("Email is required")
	}
	if password == "" {
		return errors.New("Password is required")
	}
	return nil
}

func passwordLen(password string) error {
	if len(password) < 8 {
		return errors.New("Password is short")
	}
	return nil
}

func emailValidate(email string) error {
	if !strings.Contains(email, "@") {
		return errors.New("Email is not exist")
	}
	return nil
}
