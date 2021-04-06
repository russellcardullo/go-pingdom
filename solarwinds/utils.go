package solarwinds

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"time"
)

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// Note there will be a new line character at the end of the output.
func ToJsonNoEscape(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

// RandString returns random string at specified length.
func RandString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Convert is typically used to convert map[string]interface{} to a struct in the same json structure.
func Convert(from interface{}, to interface{}) error {
	b, err := ToJsonNoEscape(from)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, to)
}
