package solarwinds

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// retrieveCookie returns the cookie value by name from http.Response.
func retrieveCookie(resp *http.Response, name string) (string, error) {
	if cookies := resp.Header[headerNameSetCookie]; cookies != nil {
		for _, cookie := range cookies {
			if strings.HasPrefix(cookie, name+"=") {
				end := strings.Index(cookie, ";")
				pair := cookie[:end]
				return strings.Split(pair, "=")[1], nil
			}
		}
		return "", fmt.Errorf("cookie '%v' does not exist in the response", name)
	} else {
		return "", errors.New("there is no cookie in the response")
	}
}
