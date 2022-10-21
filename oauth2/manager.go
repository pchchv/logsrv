package oauth2

import (
	"fmt"
	"net/http"
)

// Responsible for handling user requests in the oauth flow
// Must pick the right configuration and run the oauth redirect
type Manager struct {
	configs      map[string]Config
	startFlow    func(cfg Config, w http.ResponseWriter) error
	authenticate func(cfg Config, r *http.Request) (TokenInfo, error)
}

// NewManager creates a new Manager
func NewManager() *Manager {
	return &Manager{
		configs:      map[string]Config{},
		startFlow:    StartFlow,
		authenticate: Authenticate,
	}
}

// AddConfig for a provider
func (manager *Manager) AddConfig(providerName string, opts map[string]string) error {
	p, exist := GetProvider(providerName)

	if !exist {
		return fmt.Errorf("no provider for name %v", providerName)
	}

	cfg := Config{
		Provider: p,
		AuthURL:  p.AuthURL,
		TokenURL: p.TokenURL,
	}

	clientID, exist := opts["client_id"]
	if !exist {
		return fmt.Errorf("missing parameter client_id")
	}
	cfg.ClientID = clientID

	clientSecret, exist := opts["client_secret"]
	if !exist {
		return fmt.Errorf("missing parameter client_secret")
	}
	cfg.ClientSecret = clientSecret

	if scope, exist := opts["scope"]; exist {
		cfg.Scope = scope
	} else {
		cfg.Scope = p.DefaultScopes
	}

	if redirectURI, exist := opts["redirect_uri"]; exist {
		cfg.RedirectURI = redirectURI
	}

	manager.configs[providerName] = cfg
	return nil
}

// GetConfigs of the manager
func (manager *Manager) GetConfigs() map[string]Config {
	return manager.configs
}
