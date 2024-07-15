package util

import (
	"net/mail"
	"net/url"
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

func ValidURL(s string) bool {
	parsedURL, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return false
	}
	return true
}
