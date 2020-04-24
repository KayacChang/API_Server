package utils

import (
	"api/model"
	"fmt"
	"net/http"
	"strings"

	"github.com/badoux/checkmail"
)

func CheckMail(email string) error {

	err := checkmail.ValidateFormat(email)

	if smtpErr, ok := err.(checkmail.SmtpError); ok && err != nil {

		return &model.Error{
			Code:    http.StatusBadRequest,
			Message: smtpErr.Error(),
		}
	}

	return nil
}

func CheckSession(session string) *model.Error {

	if session == "" {
		return &model.Error{

			Code:    http.StatusBadRequest,
			Message: "Request header must contain session",
		}
	}

	return nil
}

func CheckContentType(contentType string, compare string) *model.Error {

	if contentType != compare {

		return &model.Error{
			Code:    http.StatusBadRequest,
			Message: "Content-Type must be " + compare,
		}
	}

	return nil
}

func CheckPayload(json map[string]interface{}, keys ...string) error {

	for _, key := range keys {

		if _, exist := json[key]; !exist {

			return &model.Error{
				Code:    http.StatusBadRequest,
				Message: "Request payload must contain " + key,
			}
		}
	}

	return nil
}

func CheckPassword(password string) error {

	main := func() error {
		min, max := 8, 30
		length := len(password)

		if !inRange(min, max, length) {
			return fmt.Errorf("Password must be between %d and %d characters", min, max)
		}

		hasUpperCase := false
		hasLowerCase := false
		hasNumber := false
		hasSpecial := false

		for _, char := range password {

			if isUpperCase(char) {
				hasUpperCase = true
			}

			if isLowerCase(char) {
				hasLowerCase = true
			}

			if isNumber(char) {
				hasNumber = true
			}

			if isSpecialChar(char) {
				hasSpecial = true
			}
		}

		if !hasUpperCase {
			return fmt.Errorf("Password must contain at least one uppercase")
		}

		if !hasLowerCase {
			return fmt.Errorf("Password must contain at least one lowercase letter")
		}

		if !hasNumber {
			return fmt.Errorf("Password must contain at least one number digit")
		}

		if !hasSpecial {
			return fmt.Errorf("Password must contain at least one special character")
		}

		return nil
	}

	if err := main(); err != nil {

		return &model.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	return nil
}

func inRange(min, max, tar int) bool {
	return tar >= min || tar <= max
}

func isUpperCase(char rune) bool {
	return inRange(65, 90, int(char))
}

func isLowerCase(char rune) bool {
	return inRange(97, 122, int(char))
}

func isNumber(char rune) bool {
	return inRange(48, 57, int(char))
}

func isSpecialChar(char rune) bool {
	return strings.ContainsRune(`\^$.|?*+-[]{}()`, char)
}
