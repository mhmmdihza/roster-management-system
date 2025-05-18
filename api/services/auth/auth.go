package auth

import (
	"context"
	"errors"

	st "payd/storage"

	kratos "github.com/ory/kratos-client-go"
)

var ErrAlreadyExists = errors.New("already exists")
var ErrInvalidEmail = errors.New("invalid email")
var ErrInvalidPassword = errors.New("invalid password")
var ErrInvalidCredential = errors.New("invalid credential")
var ErrTraitsKeyNotFound = errors.New("key not found")
var ErrTraitsInvalidType = errors.New("invalid type")
var ErrNotFound = errors.New("not found")
var ErrNotYetActivatingAccount = errors.New("the user has not yet activated the account")

type storage interface {
	CreateNewEmployee(ctx context.Context, name string, status string, roleId int) (int, error)
	SelectEmployeeByID(ctx context.Context, id int) (*st.Employee, error)

	NewTransacton(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Auth struct {
	kratosAdmin  *kratos.APIClient
	kratosPublic *kratos.APIClient
	storage      storage
}

type AuthOption func(*Auth) error

func NewAuth(storage storage, opts ...AuthOption) (*Auth, error) {
	auth := &Auth{storage: storage}
	for _, opt := range opts {
		if err := opt(auth); err != nil {
			return nil, err
		}
	}

	return auth, nil
}

func (a *Auth) BootstrapAdminAccount(email, name, password string) error {
	ctx := context.Background()
	id, err := a.RegisterNewUser(ctx, email, 0, true)
	if err == ErrAlreadyExists {
		return nil
	}
	if err := a.ActivateNewUser(ctx, id, name, password); err != nil {
		return err
	}
	return nil
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
