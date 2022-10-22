package login

import (
	"errors"

	"github.com/pchchv/logsrv/model"
)

const SimpleProviderName = "simple"

func init() {
	RegisterProvider(
		&ProviderDescription{
			Name:     SimpleProviderName,
			HelpText: "Simple login backend opts: user1=password,user2=password,..",
		},
		SimpleBackendFactory)
}

// Returns a new configured SimpleBackend
func SimpleBackendFactory(config map[string]string) (Backend, error) {
	userPassword := map[string]string{}
	for k, v := range config {
		userPassword[k] = v
	}
	if len(userPassword) == 0 {
		return nil, errors.New("no users provided for simple backend")
	}
	return NewSimpleBackend(userPassword), nil
}

// Working on a map of username password pairs
type SimpleBackend struct {
	userPassword map[string]string
}

// Creates a new SIMPLE Backend and verifies the parameters
func NewSimpleBackend(userPassword map[string]string) *SimpleBackend {
	return &SimpleBackend{
		userPassword: userPassword,
	}
}

// Authenticate the user
func (sb *SimpleBackend) Authenticate(username, password string) (bool, model.UserInfo, error) {
	if p, exist := sb.userPassword[username]; exist && p == password {
		return true, model.UserInfo{
			Origin: SimpleProviderName,
			Sub:    username,
		}, nil
	}
	return false, model.UserInfo{}, nil
}