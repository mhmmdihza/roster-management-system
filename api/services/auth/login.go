package auth

import (
	"context"
	"fmt"
	"net/http"
	"payd/util"
	"strconv"

	kratos "github.com/ory/kratos-client-go"
)

func (a *Auth) Login(ctx context.Context, username, password string) (*Identity, error) {
	resLoginFlow, _, err := a.kratosPublic.FrontendAPI.
		CreateNativeLoginFlow(ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("unexpected kratos login error %v", err)
	}
	flowId := resLoginFlow.Id

	res, httpResp, err := a.kratosPublic.FrontendAPI.
		UpdateLoginFlow(ctx).
		Flow(flowId).
		UpdateLoginFlowBody(kratos.UpdateLoginFlowBody{
			UpdateLoginFlowWithPasswordMethod: &kratos.UpdateLoginFlowWithPasswordMethod{
				Identifier: username,
				Password:   password,
				Method:     "password",
			},
		}).Execute()

	if err != nil || httpResp.StatusCode != http.StatusOK {
		switch httpResp.StatusCode {
		case 400:
			return nil, ErrInvalidCredential
		}
		util.Log().WithContext(ctx).WithError(err).Error("unhandled error")
		return nil, err
	}

	identity := res.GetSession().Identity
	if identity == nil {
		return nil, fmt.Errorf("unexpected identity nil")
	}
	if identity.State == nil || *identity.State != "active" {
		return nil, ErrNotYetActivatingAccount
	}
	traits, ok := identity.Traits.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected traits nil")
	}
	employeeIdStr, err := castTrait[string](traits, "employee_id")
	if err != nil {
		return nil, fmt.Errorf("can't fetch employee_id from traits: %v", err)
	}
	employeeId, err := strconv.Atoi(employeeIdStr)
	if err != nil || employeeId == 0 {
		return nil, fmt.Errorf("unexpected employee_id : %v", employeeIdStr)
	}

	employee, err := a.storage.SelectEmployeeByID(ctx, employeeId)
	if err != nil {
		return nil, err
	}
	return a.newIdentityStruct(identity.Id, traits, employee), nil
}
