package validation

import (
	"fmt"
	"net/mail"
	"regexp"
)

const (
	NAME_MIN_LENGTH = 3
	NAME_MAX_LENGTH = 100

	PASSWORD_MIN_LENGTH = 6
	PASSWORD_MAX_LENGTH = 40

	EMAIL_MIN_LENGTH = 6
	EMAIL_MAX_LENGTH = 200
)

var (
	isValidateUserName = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidateFullName = regexp.MustCompile(`^[a-zA-Z0-9_\\s]+$`).MatchString
)

func ValidateName(value string) error {
	if err := ValidateString(value, NAME_MIN_LENGTH, NAME_MAX_LENGTH); err != nil {
		return err
	}	

	if ok :=isValidateUserName(value); !ok {
		return fmt.Errorf("must contain only letter, digits, or undercore")
	}

	return nil
}

func ValidateFullName(value string) error {
	if err := ValidateString(value, NAME_MIN_LENGTH, NAME_MAX_LENGTH); err != nil {
		return err
	}

	if ok :=isValidateFullName(value); !ok {
		return fmt.Errorf("must contain only letter, space, digits, or undercore")
	}

	return nil
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, EMAIL_MIN_LENGTH, EMAIL_MAX_LENGTH); err != nil {
		return err
	}
	
	if _, err := mail.ParseAddress(value); err != nil {
		return err
	}

	return nil
}

func ValidatePassword(value string) error {
	if err := ValidateString(value, PASSWORD_MIN_LENGTH, PASSWORD_MAX_LENGTH); err != nil {
		return err
	}

	return nil
}