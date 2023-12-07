package utils

import (
	"net/http"
	"strings"
	"time"
)

// Removes whitespaces from both ends of a string and
// returns its length
func GetStringAbsoluteLength(str string) int {
	return len(strings.TrimSpace(str))
}

func GenerateCookie(name string, value string, domain string, cookiePath string, httpOnly bool, expires time.Time) http.Cookie {
	c := http.Cookie{
		Name:     name,
		Value:    value,
		Domain:   domain,
		Path:     cookiePath,
		HttpOnly: httpOnly,
		Expires:  expires,
	}

	return c
}
