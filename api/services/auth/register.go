package auth

import (
	"context"
	"payd/util"

	kratos "github.com/ory/kratos-client-go"
)

var inactiveState = "inactive"
var activeState = "active"

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
