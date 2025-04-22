package util

import (
	"fmt"
	"net/mail"
	"regexp"
	"time"
)

// ValidateString validates if the string is not empty and its length is between min and max
func ValidateString(value string, minLength, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}
	return nil
}

// ValidateUsername validates if the username meets the requirements
func ValidateUsername(username string) error {
	if err := ValidateString(username, 3, 50); err != nil {
		return err
	}

	// Check if username contains only allowed characters
	matched, err := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
	if err != nil {
		return fmt.Errorf("error validating username: %w", err)
	}
	if !matched {
		return fmt.Errorf("username can only contain letters, numbers, and underscores")
	}

	return nil
}

// ValidateEmail validates if the email is in the correct format
func ValidateEmail(email string) error {
	if err := ValidateString(email, 3, 100); err != nil {
		return err
	}

	// Check if email is in valid format
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email format: %w", err)
	}

	return nil
}

// ValidatePassword validates if the password meets the requirements
func ValidatePassword(password string) error {
	if err := ValidateString(password, 6, 100); err != nil {
		return err
	}

	// Check if password contains at least one character from each category
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z\d]`).MatchString(password)

	if !hasLower || !hasUpper || !hasDigit || !hasSpecial {
		return fmt.Errorf("password must contain at least one lowercase letter, uppercase letter, digit, and special character")
	}

	return nil
}

// ValidateDate validates if the date is not before min or after max
func ValidateDate(date time.Time, min, max time.Time) error {
	if date.Before(min) {
		return fmt.Errorf("date cannot be before %s", min.Format("2006-01-02"))
	}

	if !max.IsZero() && date.After(max) {
		return fmt.Errorf("date cannot be after %s", max.Format("2006-01-02"))
	}

	return nil
}

// ValidateURL validates if the URL is in the correct format
func ValidateURL(url string) error {
	if err := ValidateString(url, 3, 255); err != nil {
		return err
	}

	// Basic URL validation - this is a simple check
	matched, err := regexp.MatchString(`^(http|https)://[a-zA-Z0-9][-a-zA-Z0-9.]*`, url)
	if err != nil {
		return fmt.Errorf("error validating URL: %w", err)
	}
	if !matched {
		return fmt.Errorf("invalid URL format")
	}

	return nil
}

// ValidatePhoneNumber validates if the phone number is in the correct format
func ValidatePhoneNumber(phone string) error {
	if err := ValidateString(phone, 7, 20); err != nil {
		return err
	}

	// Basic phone number validation - this is a simple check
	matched, err := regexp.MatchString(`^[+]?[\d\s()-]+$`, phone)
	if err != nil {
		return fmt.Errorf("error validating phone number: %w", err)
	}
	if !matched {
		return fmt.Errorf("invalid phone number format")
	}

	return nil
}
