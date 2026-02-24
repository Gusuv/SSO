package validation

import "errors"

func RegisterValidation(username, email, password string) error {
	if err := requiredRegster(username, email, password); err != nil {
		return err
	}
	if err := passwordLen(password); err != nil {
		return err
	}
	// In development
	return nil
}

func requiredRegster(username, email, password string) error {

	if username == "" {
		return errors.New("Username is require")
	}
	if email == "" {
		return errors.New("Email is required")
	}
	if password == "" {
		return errors.New("Password is required")
	}
	return nil
}
