package auth

import (
	"errors"

	kratos "github.com/ory/kratos-client-go"
)

var ErrAlreadyExists = errors.New("already exists")
var ErrInvalidEmail = errors.New("invalid email")

type Auth struct {
	kratosAdmin  *kratos.APIClient
	kratosPublic *kratos.APIClient
}

type AuthOption func(*Auth) error

func NewAuth(opts ...AuthOption) (*Auth, error) {
	auth := &Auth{}
	for _, opt := range opts {
		if err := opt(auth); err != nil {
			return nil, err
		}
	}
	return auth, nil
}

func WithKratosPublicURL(url string) AuthOption {
	return func(a *Auth) error {
		a.kratosPublic = initKratos(url)
		return nil
	}
}

func WithKratosAdminURL(url string) AuthOption {
	return func(a *Auth) error {
		a.kratosAdmin = initKratos(url)
		return nil
	}
}

func initKratos(url string) *kratos.APIClient {
	config := kratos.NewConfiguration()
	config.Servers = []kratos.ServerConfiguration{
		{
			URL: url,
		},
	}
	return kratos.NewAPIClient(config)
}
