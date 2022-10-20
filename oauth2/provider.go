package oauth2

import "github.com/pchchv/logsrv/model"

// Description of the oauth provider adapter
type Provider struct {
	// Access the provider in the configuration
	Name string
	// Oauth authentication url to redirect to
	AuthURL string
	// Url for token exchange
	TokenURL string
	// Space separated list of oauth scopes to use for this provider
	// This list can be overwritten by configuration
	DefaultScopes string
	// Provider specific Implementation for fetching the user information
	// Possible keys in the returned map are: username, email, name
	GetUserInfo func(token TokenInfo) (u model.UserInfo, rawUserJson string, err error)
}

var provider = map[string]Provider{}

// Returns the names of all registered provider
func ProviderList() []string {
	list := make([]string, 0, len(provider))
	for k := range provider {
		list = append(list, k)
	}
	return list
}
