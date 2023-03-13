package validation

import "fmt"

func ValidateString(value string, minLength int, maxLength int) error {
	len := len(value)
	if len > maxLength || len < minLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}

	return nil
}