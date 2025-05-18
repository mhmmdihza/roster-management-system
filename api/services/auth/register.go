package auth

import (
	"context"
	"fmt"
	"net/http"
	"payd/util"
	"strconv"

	kratos "github.com/ory/kratos-client-go"
)

var inactiveState = "inactive"
var activeState = "active"

type Identity struct {
	ID          string // kratos userid
	Email       string // registered email
	EmployeeId  string // db employee id
	Role        string // privilege-based(admin/employee)
	PrimaryRole int    // responsibility-based
}

// kratos traits schema
func (i *Identity) GetTraits() map[string]interface{} {
	return map[string]interface{}{
		"email":        i.Email,
		"role":         i.Role,
		"primary_role": i.PrimaryRole,
		"employee_id":  i.EmployeeId,
	}
}

// register a new user in an inactive state
// later, the user will activate the account and fill in the password and other information
// roleAdmin = privilege-based meaning refers to the level of access, permissions, or authority
// primaryRole = responsibility-based meaning refers to the main job or task someone is assigned to do
func (a *Auth) RegisterNewUser(ctx context.Context, email string, primaryRole int, roleAdmin bool) (string, error) {
	role := "employee"
	if roleAdmin {
		role = "admin"
	}
	traits := map[string]interface{}{
		"email":        email,
		"role":         role,
		"primary_role": primaryRole,
	}
	pass := "123456"
	identity, httpResp, err := a.kratosAdmin.IdentityAPI.CreateIdentity(ctx).
		CreateIdentityBody(kratos.CreateIdentityBody{
			Credentials: &kratos.IdentityWithCredentials{
				Password: &kratos.IdentityWithCredentialsPassword{
					Config: &kratos.IdentityWithCredentialsPasswordConfig{
						Password: &pass,
					},
				},
			},
			SchemaId: "default",
			Traits:   traits,
			State:    &inactiveState,
		}).Execute()
	if err != nil {
		if httpResp == nil {
			return "", err
		}
		switch httpResp.StatusCode {
		case 409:
			return "", ErrAlreadyExists
		case 400:
			return "", ErrInvalidEmail
		}

		util.Log().WithContext(ctx).WithError(err).Error("unhandled error")
		return "", err
	}
	return identity.Id, nil
}

func (a *Auth) GetIdentity(ctx context.Context, userId string) (*Identity, error) {
	identity, httpResp, err := a.kratosAdmin.IdentityAPI.GetIdentity(ctx, userId).Execute()
	if err != nil {
		if httpResp == nil {
			return nil, err
		}
		switch httpResp.StatusCode {
		case 404:
			return nil, ErrNotFound
		}
		util.Log().WithContext(ctx).WithError(err).Error("unhandled error")
		return nil, err
	}
	if identity != nil && identity.State != nil && *identity.State == activeState {
		return nil, ErrAlreadyExists
	}
	traits, ok := identity.Traits.(map[string]interface{})
	if !ok {
		// should be considered as internal error
		return nil, fmt.Errorf("expecting map traits")
	}
	primaryRole, err := castTrait[float64](traits, "primary_role")
	if err != nil {
		return nil, err
	}
	email, err := castTrait[string](traits, "email")
	if err != nil {
		return nil, err
	}
	role, err := castTrait[string](traits, "role")
	if err != nil {
		return nil, err
	}
	return &Identity{
		ID:          userId,
		Email:       email,
		PrimaryRole: int(primaryRole),
		Role:        role,
	}, nil
}

func (a *Auth) ActivateNewUser(ctx context.Context, userId, name, password string) error {
	identity, err := a.GetIdentity(ctx, userId)
	if err != nil {
		return err
	}
	tctx, err := a.storage.NewTransacton(ctx)
	if err != nil {
		util.Log().WithContext(ctx).WithError(err).Error("failed to start db transaction")
		return err
	}
	defer a.dbTransactions(tctx, err)
	employeeId, err := a.storage.CreateNewEmployee(tctx, name, "ACTIVE", int(identity.PrimaryRole))
	if err != nil {
		util.Log().WithContext(tctx).WithError(err).Error("storage create new employee")
		return err
	}
	traits := identity.GetTraits()
	traits["employee_id"] = strconv.Itoa(employeeId)

	_, httpResp, err := a.kratosAdmin.IdentityAPI.UpdateIdentity(tctx, userId).
		UpdateIdentityBody(kratos.UpdateIdentityBody{
			Credentials: &kratos.IdentityWithCredentials{
				Password: &kratos.IdentityWithCredentialsPassword{
					Config: &kratos.IdentityWithCredentialsPasswordConfig{
						Password: &password,
					},
				},
			},
			SchemaId: "default",
			Traits:   traits,
			State:    activeState,
		}).Execute()
	if err != nil || httpResp.StatusCode != http.StatusOK {
		switch httpResp.StatusCode {
		case 400:
			return ErrInvalidPassword
		}
		util.Log().WithContext(ctx).WithError(err).Error("unhandled error")
		return err
	}
	return nil
}

// defer only after calling storage.NewTransacton
func (a *Auth) dbTransactions(ctx context.Context, err error) {
	if err != nil {
		util.Log().WithContext(ctx).WithError(err).Error("rollback db trx")
		err = a.storage.Rollback(ctx)
		if err != nil {
			util.Log().WithContext(ctx).WithError(err).Error("failed rollback")
		}
		return
	}
	err = a.storage.Commit(ctx)
	if err != nil {
		util.Log().WithContext(ctx).WithError(err).Error("failed commit")
	}
}

func castTrait[T any](m map[string]interface{}, key string) (T, error) {
	val, exists := m[key]
	if !exists {
		var zero T
		return zero, ErrTraitsKeyNotFound
	}

	casted, ok := val.(T)
	if !ok {
		var zero T
		return zero, ErrTraitsInvalidType
	}

	return casted, nil
}
