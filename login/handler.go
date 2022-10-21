package login

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/pchchv/logsrv/model"
	"github.com/pchchv/logsrv/oauth2"
)

// Mail login handler.
// It serves the login ressource and does the authentication against the backends or oauth provider.
type Handler struct {
	backends         []Backend
	oauth            oauthManager
	config           *Config
	signingMethod    jwt.SigningMethod
	signingKey       interface{}
	signingVerifyKey interface{}
	userClaims       userClaimsFunc
}

type userClaimsFunc func(userInfo model.UserInfo) (jwt.Claims, error)

type oauthManager interface {
	Handle(w http.ResponseWriter, r *http.Request) (
		startedFlow bool,
		authenticated bool,
		userInfo model.UserInfo,
		err error)
	AddConfig(providerName string, opts map[string]string) error
	GetConfigFromRequest(r *http.Request) (oauth2.Config, error)
}

// Creates a login handler based on the supplied configuration
func NewHandler(config *Config) (*Handler, error) {
	if len(config.Backends) == 0 && len(config.Oauth) == 0 {
		return nil, errors.New("No login backends or oauth provider configured")
	}

	backends := []Backend{}
	for pName, opts := range config.Backends {
		p, exist := GetProvider(pName)
		if !exist {
			return nil, fmt.Errorf("No such provider: %v", pName)
		}
		b, err := p(opts)
		if err != nil {
			return nil, err
		}
		backends = append(backends, b)
	}

	oauth := oauth2.NewManager()
	for providerName, opts := range config.Oauth {
		err := oauth.AddConfig(providerName, opts)
		if err != nil {
			return nil, err
		}
	}

	userClaims, err := NewUserClaims(config)
	if err != nil {
		return nil, err
	}

	return &Handler{
		backends:   backends,
		config:     config,
		oauth:      oauth,
		userClaims: userClaims.Claims,
	}, nil
}
