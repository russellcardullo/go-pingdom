package solarwinds

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRetrieveCookie(t *testing.T) {
	setup()
	defer teardown()

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(headerNameSetCookie, "Swicus-auth=; Path=/; Expires=Thu, 01 Jan 1970 00:00:00 GMT; HttpOnly; Secure")
		w.Header().Add(headerNameSetCookie, "swicus=abcd; Path=/; Expires=Tue, 06 Apr 2021 09:18:16 GMT; Max-Age=1209600; HttpOnly; Secure; SameSite=None")
	}
	req := httptest.NewRequest("GET", "http://foo.com", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	resp := w.Result()
	swicus, err := retrieveCookie(resp, "swicus")
	assert.NoError(t, err)
	assert.Equal(t, "abcd", swicus)
}
