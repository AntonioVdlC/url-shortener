package utils

import (
	"log"
	"net/url"
)

// TODO: maybe also add https://developers.google.com/safe-browsing/v4
func IsValidUrl(str string) (bool, string) {
	log.Printf("Validating '%s'", str)

	u, err := url.Parse(str)

	if err != nil {
		return false, "Invalid url format"
	}

	if u.Scheme == "" {
		return false, "Missing scheme"
	}

	if u.Scheme != "https" {
		return false, "Scheme not secure"
	}

	if u.Host == "" {
		return false, "Missing host"
	}

	return true, ""
}
