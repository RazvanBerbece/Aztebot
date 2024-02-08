package rolesService

import (
	dataModels "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/models"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
)

// Returns the highest staff and order roles from an array of roles, or nil if not applicable
func GetHighestRoles(roles []dataModels.Role) (*dataModels.Role, *dataModels.Role) {

	var highestStaffIdx = -1
	var highestOrderIdx = -1

	for idx, role := range roles {
		if utils.RoleIsStaffRole(role.Id) && idx > highestStaffIdx {
			highestStaffIdx = idx
		} else if !utils.RoleIsStaffRole(role.Id) && idx > highestOrderIdx {
			highestOrderIdx = idx
		}
	}

	return &roles[highestStaffIdx], &roles[highestOrderIdx]

}
