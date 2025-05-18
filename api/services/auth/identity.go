package auth

import (
	"errors"
	st "payd/storage"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Identity struct {
	ID           string // kratos userid
	Email        string // registered email
	EmployeeId   string // db employee id
	EmployeeName string // db employee name
	Role         string // privilege-based(admin/employee)
	PrimaryRole  int    // responsibility-based from db role_id

	jwtSecret []byte
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

func (i *Identity) GenerateJWT(expiration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":           i.ID,
		"email":         i.Email,
		"employee_id":   i.EmployeeId,
		"employee_name": i.EmployeeName,
		"role":          i.Role,
		"primary_role":  i.PrimaryRole,
		"exp":           time.Now().Add(expiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(i.jwtSecret)
}

func (a *Auth) VerifySignatureJWT(tokenStr string) (*Identity, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return a.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims type")
	}

	exp := int64(claims["exp"].(float64))
	if time.Now().Unix() > exp {
		return nil, errors.New("token expired")
	}

	return &Identity{
		ID:           claims["sub"].(string),
		Email:        claims["email"].(string),
		EmployeeId:   claims["employee_id"].(string),
		EmployeeName: claims["employee_name"].(string),
		Role:         claims["role"].(string),
		PrimaryRole:  int(claims["primary_role"].(float64)),
	}, nil
}

// build identity struct based on kratos traits map and employee table db and set jwt secret
func (a *Auth) newIdentityStruct(kratosId string, traits map[string]interface{}, employee *st.Employee) *Identity {
	identity := &Identity{jwtSecret: a.jwtSecret}
	if employee != nil {
		identity.EmployeeId = strconv.Itoa(employee.ID)
		identity.EmployeeName = employee.Name
		identity.PrimaryRole = employee.PrimaryRole
	}
	identity.ID = kratosId
	email, err := castTrait[string](traits, "email")
	if err == nil {
		identity.Email = email
	}
	role, err := castTrait[string](traits, "role")
	if err == nil {
		identity.Role = role
	}
	return identity
}
