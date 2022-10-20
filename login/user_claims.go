package login

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pchchv/logsrv/model"
	"github.com/pkg/errors"
)

type customClaims map[string]interface{}

type UserClaims interface {
	Claims(userInfo model.UserInfo) (jwt.Claims, error)
}

func (custom customClaims) Valid() error {
	if exp, ok := custom["exp"]; ok {
		if exp, ok := exp.(int64); ok {
			if exp < time.Now().Unix() {
				return errors.New("token expired")
			}
		}
	}
	return nil
}

func (custom customClaims) merge(values map[string]interface{}) {
	for k, v := range values {
		custom[k] = v
	}
}

func NewUserClaims(config *Config) (UserClaims, error) {
	if config.UserEndpoint != "" {
		return newUserClaimsProvider(config.UserEndpoint, config.UserEndpointToken, config.UserEndpointTimeout)
	}
	return newUserClaimsFile(config.UserFile)
}
