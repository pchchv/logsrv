package oauth2

import "net/http"

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
