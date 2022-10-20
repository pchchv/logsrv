package login

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

// Config for the loginsrv handler
type Config struct {
	Host                   string
	Port                   string
	LogLevel               string
	TextLogging            bool
	JwtSecret              string
	JwtSecretFile          string
	JwtAlgo                string
	JwtExpiry              time.Duration
	JwtRefreshes           int
	SuccessURL             string
	Redirect               bool
	RedirectQueryParameter string
	RedirectCheckReferer   bool
	RedirectHostFile       string
	LogoutURL              string
	Template               string
	LoginPath              string
	CookieName             string
	CookieExpiry           time.Duration
	CookieDomain           string
	CookieHTTPOnly         bool
	CookieSecure           bool
	Backends               Options
	Oauth                  Options
	GracePeriod            time.Duration
	UserFile               string
	UserEndpoint           string
	UserEndpointToken      string
	UserEndpointTimeout    time.Duration
}

// Configuration structure for oauth and backend provider
// key is the providername, value is a options map
type Options map[string]map[string]string

var jwtDefaultSecret string

func init() {
	var err error
	jwtDefaultSecret, err = randStringBytes(64)
	if err != nil {
		panic(err)
	}
}

// Adds the options for a provider in the form of key=value,key=value...
func (c *Config) addOauthOpts(providerName, optsKvList string) error {
	opts, err := parseOptions(optsKvList)
	if err != nil {
		return err
	}
	c.Oauth[providerName] = opts
	return nil
}

// Adds the options for a provider in the form of key=value,key=value...
func (c *Config) addBackendOpts(providerName, optsKvList string) error {
	opts, err := parseOptions(optsKvList)
	if err != nil {
		return err
	}
	c.Backends[providerName] = opts
	return nil
}

// Resolves configuration values, which are dynamically referenced via files
func (c *Config) ResolveFileReferences() error {
	// Try to load the secret from a file, if set
	if c.JwtSecretFile != "" {
		secretBytes, err := ioutil.ReadFile(c.JwtSecretFile)
		if err != nil {
			return err
		}
		c.JwtSecret = string(secretBytes)
	}
	return nil
}

// Default config for the loginsrv handler
func DefaultConfig() *Config {
	return &Config{
		Host:                   "localhost",
		Port:                   "6789",
		LogLevel:               "info",
		JwtSecret:              jwtDefaultSecret,
		JwtAlgo:                "HS512",
		JwtExpiry:              24 * time.Hour,
		JwtRefreshes:           0,
		SuccessURL:             "/",
		Redirect:               true,
		RedirectQueryParameter: "backTo",
		RedirectCheckReferer:   true,
		RedirectHostFile:       "",
		LogoutURL:              "",
		LoginPath:              "/login",
		CookieName:             "jwt_token",
		CookieHTTPOnly:         true,
		CookieSecure:           true,
		Backends:               Options{},
		Oauth:                  Options{},
		GracePeriod:            5 * time.Second,
		UserFile:               "",
		UserEndpoint:           "",
		UserEndpointToken:      "",
		UserEndpointTimeout:    5 * time.Second,
	}
}

func randStringBytes(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func parseOptions(b string) (map[string]string, error) {
	opts := map[string]string{}
	pairs := strings.Split(b, ",")
	for _, p := range pairs {
		pair := strings.SplitN(p, "=", 2)
		if len(pair) != 2 {
			return nil, fmt.Errorf("provider configuration has to be in form 'key1=value1,key2=..', but was %v", p)
		}
		opts[pair[0]] = pair[1]
	}
	return opts, nil
}
