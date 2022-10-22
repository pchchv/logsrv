package osiam

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/pchchv/logsrv/logging"
	"github.com/pchchv/logsrv/login"
	"github.com/pchchv/logsrv/model"
)

// OSIAM authentication backend
type Backend struct {
	client *Client
}

const OsiamProviderName = "osiam"

func init() {
	login.RegisterProvider(
		&login.ProviderDescription{
			Name:     OsiamProviderName,
			HelpText: "Osiam login backend opts: endpoint=..,client_id=..,client_secret=..",
		},
		func(config map[string]string) (login.Backend, error) {
			if config["clientId"] != "" {
				logging.Logger.Warn("DEPRECATED: please use 'client_id' and 'client_secret' in future.")
				return NewBackend(config["endpoint"], config["clientId"], config["clientSecret"])
			}
			return NewBackend(config["endpoint"], config["client_id"], config["client_secret"])
		})
}

// Creates a new OSIAM Backend and verifies the parameters
func NewBackend(endpoint, clientID, clientSecret string) (*Backend, error) {
	if _, err := url.Parse(endpoint); err != nil {
		return nil, fmt.Errorf("osiam endpoint has to be a valid url: %v: %v", endpoint, err)
	}
	if clientID == "" {
		return nil, errors.New("no osiam clientID provided.")
	}
	if clientSecret == "" {
		return nil, errors.New("no osiam clientSecret provided")
	}
	client := NewClient(endpoint, clientID, clientSecret)
	return &Backend{
		client: client,
	}, nil
}

// Authenticate the user
func (b *Backend) Authenticate(username, password string) (bool, model.UserInfo, error) {
	authenticated, _, err := b.client.GetTokenByPassword(username, password)
	if !authenticated || err != nil {
		return authenticated, model.UserInfo{}, err
	}
	userInfo := model.UserInfo{
		Origin: OsiamProviderName,
		Sub:    username,
	}
	return true, userInfo, nil
}
