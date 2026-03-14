package validation

import "regexp"

func LoginValidation(email, password string, appId int64) error {
	if err := emailValidate(email); err != nil {
		return err
	}
	if err := passwordValidate(password); err != nil {
		return err
	}
	if appId == 0 {
		return RequiredAppId
	}

	return nil
}

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-z]{2,}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
)

const (
	minPasswordLen = 8
	maxPasswordLen = 60
	minUsernameLen = 2
	maxUsernameLen = 40
)

func emailValidate(email string) error {
	if email == "" {
		return RequiredEmail
	}
	if !emailRegex.MatchString(email) {
		return InvalidEmail
	}
	return nil
}

func passwordValidate(password string) error {
	if password == "" {
		return RequiredPassword
	}
	if len(password) < minPasswordLen {
		return ShortPassword
	}
	if len(password) > maxPasswordLen {
		return LongPassword
	}
	return nil
}
