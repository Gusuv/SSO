package validation

import "errors"

func RegisterValidation(username, email, password string) error {
	if err := requiredRegister(username, email, password); err != nil {
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

func requiredRegister(username, email, password string) error {

	if username == "" {
		return errors.New("Username is required")
	}
	if email == "" {
		return errors.New("Email is required")
	}
	if password == "" {
		return errors.New("Password is required")
	}
	return nil
}
