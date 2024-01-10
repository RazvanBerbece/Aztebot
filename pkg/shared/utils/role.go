package utils

func GetCircleAndOrderForGivenRoles(roleIds []int) (string, *int) {

	var circle string
	var order *int

	var hasInnerCircleId bool = false
	var maxInnerOrderId int = -1
	for _, roleId := range roleIds {
		circle, order := GetCircleAndOrderFromRoleId(roleId)
		if circle == 1 {
			hasInnerCircleId = true
			if order > maxInnerOrderId {
				maxInnerOrderId = order
			}
		}
	}

	if hasInnerCircleId {
		circle = "INNER"
	} else {
		circle = "OUTER"
	}

	if maxInnerOrderId == -1 {
		order = nil
	} else {
		order = &maxInnerOrderId
	}

	return circle, order
}

func GetCircleAndOrderFromRoleId(roleId int) (int, int) {

	if roleId <= 7 {
		return 0, -1
	} else {
		if roleId >= 7 && roleId < 12 {
			return 1, 1
		} else if roleId >= 12 && roleId < 15 {
			return 1, 2
		} else if roleId >= 15 {
			return 1, 3
		}
	}

	return 0, -1

}
