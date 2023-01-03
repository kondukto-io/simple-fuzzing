package util

import (
	"net/mail"
	"strconv"
)

func VaildID(s string) bool {
	if v, err := strconv.Atoi(s); err != nil || v < 1 {
		return false
	}

	return true
}

func isValidEmail(s string) bool {
	_, err := mail.ParseAddress(s)
	return err == nil
}
