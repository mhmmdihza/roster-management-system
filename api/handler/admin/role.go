package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoleResponse struct {
	ID       int    `json:"id"`
	RoleName string `json:"roleName"`
}

func (a *Admin) listRole(c *gin.Context) {
	var res []RoleResponse
	listRole := a.role.GetRoles()
	for _, role := range listRole {
		res = append(res, RoleResponse{
			ID:       role.ID,
			RoleName: role.Name,
		})
	}
	c.JSON(http.StatusOK, res)
}
