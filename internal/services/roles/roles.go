package rolesService

import (
	dax "github.com/RazvanBerbece/Aztebot/internal/data/models/dax/aztebot"
	globalConfiguration "github.com/RazvanBerbece/Aztebot/internal/globals/configuration"
	"github.com/RazvanBerbece/Aztebot/pkg/shared/utils"
)

// Returns the highest staff and order roles from an array of roles, or nil if not applicable
func GetHighestRoles(roles []dax.Role) (*dax.Role, *dax.Role) {

	var highestStaffIdx = -1
	var highestOrderIdx = -1

	var highestStaffRole *dax.Role
	var highestOrderRole *dax.Role

	for idx, role := range roles {
		if utils.StringInSlice(role.DisplayName, globalConfiguration.StaffRoles) && idx > highestStaffIdx {
			highestStaffIdx = idx
			highestStaffRole = &roles[highestStaffIdx]
		} else if !utils.StringInSlice(role.DisplayName, globalConfiguration.StaffRoles) && idx > highestOrderIdx {
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
