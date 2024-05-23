package utils

import (
	"strconv"
	"strings"
)

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

func GetOrderAsString(orderId int) string {
	switch orderId {
	case 1:
		return "---- First Order ----"
	case 2:
		return "---- Second Order ----"
	case 3:
		return "---- Third Order ----"
	default:
		return "not supported"
	}
}

func GetCircleAndOrderFromRoleId(roleId int) (int, int) {

	if roleId <= 10 {
		return 0, -1
	} else {
		if roleId >= 11 && roleId < 15 {
			return 1, 1
		} else if roleId >= 15 && roleId < 18 {
			return 1, 2
		} else if roleId >= 18 {
			return 1, 3
		}
	}

	return 0, -1

}

func GetRoleIdsFromRoleString(roleIdsString string) []int {

	var roleIds []int = []int{}
	roleIdsTokens := strings.Split(roleIdsString, ",")

	for _, roleIdToken := range roleIdsTokens {
		if len(roleIdToken) > 0 {
			i, err := strconv.Atoi(roleIdToken)
			if err != nil {
				return nil
			}
			roleIds = append(roleIds, int(i))
		}
	}

	return roleIds

}
