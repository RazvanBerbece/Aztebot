package rolesService

import (
	dataModels "github.com/RazvanBerbece/Aztebot/internal/bot-service/data/models"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
)

// Returns the highest staff and order roles from an array of roles, or nil if not applicable
func GetHighestRoles(roles []dataModels.Role) (*dataModels.Role, *dataModels.Role) {

	var highestStaffIdx = -1
	var highestOrderIdx = -1

	var highestStaffRole *dataModels.Role
	var highestOrderRole *dataModels.Role

	for idx, role := range roles {
		if utils.RoleIsStaffRole(role.Id) && idx > highestStaffIdx {
			highestStaffIdx = idx
			highestStaffRole = &roles[highestStaffIdx]
		} else if !utils.RoleIsStaffRole(role.Id) && idx > highestOrderIdx {
			highestOrderIdx = idx
			highestOrderRole = &roles[highestOrderIdx]
		}
	}

	if highestStaffIdx == -1 {
		highestStaffRole = nil
	}

	if highestOrderIdx == -1 {
		highestOrderRole = nil
	}

	return highestStaffRole, highestOrderRole

}
