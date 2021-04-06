package solarwinds

import "fmt"

const (
	ErrCodeNetworkException uint32 = iota
	ErrCodeDeleteActiveUserException
)

type ClientError struct {
	StatusCode uint32 `json:"statusCode"`
	Err        error  `json:"err"`
}

func (c *ClientError) Error() string {
	return fmt.Sprintf("status: %d, err: %v", c.StatusCode, c.Err)
}

func NewNetworkError(cause error) error {
	return &ClientError{
		StatusCode: ErrCodeNetworkException,
		Err:        cause,
	}
}

func NewErrorAttemptDeleteActiveUser(user string) error {
	return &ClientError{
		StatusCode: ErrCodeDeleteActiveUserException,
		Err:        fmt.Errorf("deleting active user %v is not supported", user),
	}
}
