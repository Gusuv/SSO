package validation

func RegisterValidation(username, email, password string) error {
	if err := passwordLen(password); err != nil {
		return err
	}
	// In development
	return nil
}
