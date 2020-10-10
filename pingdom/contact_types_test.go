package pingdom

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContact_ValidContact_Positive(t *testing.T) {
	name := "testName"
	contact := Contact{
		Name: name,
	}

	err := contact.ValidContact()

	assert.Equal(t, nil, err, "Contact.ValidContact() should return nil")
}

func TestContact_ValidContact_Negative(t *testing.T) {
	contact := Contact{
		Name: "",
	}

	want := fmt.Errorf("Invalid value for `Name`.  Must contain non-empty string")

	err := contact.ValidContact()

	assert.Equal(t, want, err, "Contact.ValidContact() should return error")
}
