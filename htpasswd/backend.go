package htpasswd

import (
	"errors"
	"strings"

	"github.com/pchchv/logsrv/login"
	"github.com/pchchv/logsrv/model"
)

// htpasswd based authentication backend
type Backend struct {
	auth *Auth
}

const ProviderName = "htpasswd"

func init() {
	login.RegisterProvider(
		&login.ProviderDescription{
			Name:     ProviderName,
			HelpText: "Htpasswd login backend opts: files=/path/to/pwdfile,/path/to/additionalfile",
		},
		BackendFactory)
}

// Creates a htpasswd backend
func BackendFactory(config map[string]string) (login.Backend, error) {
	var files []string
	if f, exist := config["file"]; exist {
		files = append(files, strings.Split(f, ";")...)
	}
	if len(files) == 0 {
		return nil, errors.New(`missing parameter "file" for htpasswd provider`)
	}
	return NewBackend(files)
}

// Creates a new Backend and verifies the parameters.
func NewBackend(filenames []string) (*Backend, error) {
	auth, err := NewAuth(filenames)
	return &Backend{
		auth,
	}, err
}

// Authenticate the user
func (sb *Backend) Authenticate(username, password string) (bool, model.UserInfo, error) {
	authenticated, err := sb.auth.Authenticate(username, password)
	if authenticated && err == nil {
		return authenticated, model.UserInfo{
			Origin: ProviderName,
			Sub:    username,
		}, err
	}
	return false, model.UserInfo{}, err
}
