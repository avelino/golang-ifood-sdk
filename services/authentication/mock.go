package authentication

import (
	"github.com/stretchr/testify/mock"
)

// AuthMock mock of auth service
type AuthMock struct {
	mock.Mock
}

// Authenticate mock of auth service
func (a *AuthMock) Authenticate(user, pass string) (c *Credentials, err error) {
	args := a.Called(user, pass)
	if res, ok := args.Get(0).(*Credentials); ok {
		return res, nil
	}
	return nil, args.Error(1)
}

// Validate mock of auth service
func (a *AuthMock) Validate() (err error) {
	args := a.Called()
	return args.Error(0)
}

// GetToken mock of auth service
func (a *AuthMock) GetToken() (token string) {
	args := a.Called()
	return args.Get(0).(string)
}
