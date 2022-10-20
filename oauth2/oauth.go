package oauth2

// Describes a typical 3-legged OAuth2 flow,
// with both the client application information and the server's endpoint URLs
type Config struct {
	// Application's ID
	ClientID string
	// Application's secret
	ClientSecret string
	// OAuth authentication url to redirect to
	AuthURL string
	// Url for token exchange
	TokenURL string
	// URL to redirect users going through the OAuth flow, after the resource owner's URLs
	RedirectURI string
	// Specifies optional requested permissions, this is a *space* separated list
	Scope string
	// The OAuth provider
	Provider Provider
}

// Represents the credentials used to authorize
// the requests to access protected resources on the OAuth 2.0
// provider's backend
type TokenInfo struct {
	// Token that authorizes and authenticates the requests
	AccessToken string `json:"access_token"`
	// Type of token
	TokenType string `json:"token_type,omitempty"`
	// Scopes for this tolen
	Scope string `json:"scope,omitempty"`
}
