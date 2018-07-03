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

	var want error
	want = nil

	err := user.ValidCreate()

	assert.Equal(t, want, err, "User.ValidCreate() should return nil")

}

func TestUserValidCreateNegative(t *testing.T) {
	name := ""
	user := User{
		Username : name,
	}

	want := fmt.Errorf("Invalid value for `Username`.  Must contain non-empty string")

	err := user.ValidCreate()

	assert.Equal(t, want, err, "User.ValidCreate() should return error")

}