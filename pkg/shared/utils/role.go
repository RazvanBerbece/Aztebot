package utils

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
