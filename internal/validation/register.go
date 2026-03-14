package validation

func RegisterValidation(username, email, password string) error {

	if err := usernameValidate(username); err != nil {
		return err
	}
	if err := emailValidate(email); err != nil {
		return err
	}
	if err := passwordValidate(password); err != nil {
		return err
	}
	return nil
}

func usernameValidate(username string) error {
	if username == "" {
		return RequiredUsername
	}
	if len(username) < minUsernameLen {
		return ShortUsername
	}
	if len(username) > maxUsernameLen {
		return LongUsername
	}
	if !usernameRegex.MatchString(username) {
		return InvalidUsername
	}

	return nil
}
