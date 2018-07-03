package pingdom

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"fmt"
)

func TestUserPostParams(t *testing.T) {
	name := "testUsername"

	user := User{
		Username : name,
	}
	params := user.PostParams()
	want := map[string]string{
		"name": name,
	}

	assert.Equal(t, want, params, "User.PostParams() should return correct map")
}

func TestUserValidCreatePositive(t *testing.T) {
	name := "testUsername"
	user := User{
		Username : name,
	}

	err := user.ValidCreate()

	assert.Equal(t, nil, err, "User.ValidCreate() should return nil")
}

func TestUserValidCreateNegative(t *testing.T) {
	user := User{
		Username : "",
	}

	want := fmt.Errorf("Invalid value for `Username`.  Must contain non-empty string")

	err := user.ValidCreate()

	assert.Equal(t, want, err, "User.ValidCreate() should return error")
}

func TestUserValidCreateContactPositive1(t *testing.T) {
	contact := Contact{
		Email: "test@example.com",
		CountryCode: "1",
		Number: "5559995555",
	}

	err := contact.ValidCreateContact()
	assert.Equal(t, nil, err, "contact.ValidCreateContact() should return nil")
}

func TestUserValidCreateContactPositive2(t *testing.T) {
	contact := Contact{
		Email: "test@example.com",
	}

	err := contact.ValidCreateContact()
	assert.Equal(t, nil, err, "contact.ValidCreateContact() should return nil")
}

func TestUserValidCreateContactPositive3(t *testing.T) {
	contact := Contact{
		CountryCode: "1",
		Number: "5559995555",
	}

	err := contact.ValidCreateContact()
	assert.Equal(t, nil, err, "contact.ValidCreateContact() should return nil")
}


func TestUserValidCreateContactNegative1(t *testing.T) {
	contact := Contact{
		CountryCode: "1",
	}

	want := fmt.Errorf("you must provide either an Email or a Phone Number to create a contact target")

	err := contact.ValidCreateContact()

	assert.Equal(t, want, err, "contact.ValidCreateContact() should return error")
}

func TestUserValidCreateContactNegative2(t *testing.T) {
	contact := Contact{
		Number: "5559995555",
	}

	want := fmt.Errorf("you must provide a Country Code if providing a phone number")

	err := contact.ValidCreateContact()

	assert.Equal(t, want, err, "contact.ValidCreateContact() should return error")
}

func TestPostContactParams(t *testing.T) {
	email := "test@example.com"
	countrycode := "1"
	number := "5559995555"

	contact := Contact{
		Email: email,
		CountryCode: countrycode,
		Number: number,
	}
	params := contact.PostContactParams()
	want := map[string]string{
		"email" : email,
		"number" : number,
		"countrycode" : countrycode,
	}

	assert.Equal(t, want, params, "Contact.PostContactParams() should return correct map")
}