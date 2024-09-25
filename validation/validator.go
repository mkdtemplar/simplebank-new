package validation

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUserName = regexp.MustCompile("^[a-zA-Z0-9_]+$").MatchString
	isValidFullName = regexp.MustCompile("^[a-zA-Z\\s]+$").MatchString
)

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("string length must be between %d and %d", minLength, maxLength)
	}
	return nil
}

func ValidateUserName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}

	if !isValidUserName(value) {
		return fmt.Errorf("invalid user name, only letters, numbers or underscore: %s", value)
	}

	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 6, 100)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	_, err := mail.ParseAddress(value)
	if err != nil {
		return fmt.Errorf("invalid email address: %s", value)
	}
	return nil
}

func ValidateFullName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return fmt.Errorf("length must be between %d and %d", 3, 100)
	}
	if !isValidFullName(value) {
		return fmt.Errorf("invalid full name, only letters: %s", value)
	}
	return nil
}

func VerifyEmailId(value int64) error {
	if value <= 0 {
		return fmt.Errorf("id must be positive integer")
	}
	return nil
}

func ValidateSecretCode(value string) error {
	return ValidateString(value, 32, 128)
}