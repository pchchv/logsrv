package oauth2

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

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

// Represents an oauth error response in json form
type JSONError struct {
	Error string `json:"error"`
}

const (
	stateCookieName = "oauthState"
	defaultTimeout  = 5 * time.Second
)

// Starts a flow, redirecting the user to the login provider
// A state parameter to protect against cross-site request forgery attacks is randomly generated and stored in a cookie
func StartFlow(cfg Config, w http.ResponseWriter) error {
	values := make(url.Values)
	values.Set("client_id", cfg.ClientID)
	values.Set("scope", cfg.Scope)
	values.Set("redirect_uri", cfg.RedirectURI)
	values.Set("response_type", "code")
	// Set and store the state param
	state, err := randStringBytes(32)
	if err != nil {
		return err
	}
	values.Set("state", state)
	http.SetCookie(w, &http.Cookie{
		Name:     stateCookieName,
		MaxAge:   60 * 10, // 10 minutes
		Value:    values.Get("state"),
		HttpOnly: true,
	})
	targetURL := cfg.AuthURL + "?" + values.Encode()
	w.Header().Set("Location", targetURL)
	w.WriteHeader(http.StatusFound)
	return nil
}

// Authenticate after coming back from the oauth flow
// Verify the state parameter against the state cookie from the request
func Authenticate(cfg Config, r *http.Request) (TokenInfo, error) {
	if r.FormValue("error") != "" {
		return TokenInfo{}, fmt.Errorf("error: %v", r.FormValue("error"))
	}
	state := r.FormValue("state")
	stateCookie, err := r.Cookie(stateCookieName)
	if err != nil || stateCookie.Value != state {
		return TokenInfo{}, fmt.Errorf("error: oauth state param could not be verified")
	}
	code := r.FormValue("code")
	if code == "" {
		return TokenInfo{}, fmt.Errorf("error: no auth code provided")
	}
	return getAccessToken(cfg, state, code)
}

func getAccessToken(cfg Config, state, code string) (TokenInfo, error) {
	values := url.Values{}
	values.Set("client_id", cfg.ClientID)
	values.Set("client_secret", cfg.ClientSecret)
	values.Set("code", code)
	values.Set("redirect_uri", cfg.RedirectURI)
	values.Set("grant_type", "authorization_code")
	r, _ := http.NewRequest("POST", cfg.TokenURL, strings.NewReader(values.Encode()))
	cntx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	_ = r.WithContext(cntx)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return TokenInfo{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return TokenInfo{}, fmt.Errorf("error: expected http status 200 on token exchange, but got %v", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TokenInfo{}, fmt.Errorf("error reading token exchange response: %q", err)
	}
	jsonError := JSONError{}
	err = json.Unmarshal(body, &jsonError)
	if err != nil {
		panic(err)
	}
	if jsonError.Error != "" {
		return TokenInfo{}, fmt.Errorf("error: got %q on token exchange", jsonError.Error)
	}
	tokenInfo := TokenInfo{}
	err = json.Unmarshal(body, &tokenInfo)
	if err != nil {
		return TokenInfo{}, fmt.Errorf("error on parsing oauth token: %v", err)
	}
	if tokenInfo.AccessToken == "" {
		return TokenInfo{}, fmt.Errorf("error: no access_token on token exchange")
	}
	return tokenInfo, nil
}

func randStringBytes(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
