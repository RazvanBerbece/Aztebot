package dataModels

type JailedUser struct {
	UserId            string
	Reason            string
	TaskToComplete    string
	JailedAt          int64
	RoleIdsBeforeJail string
}
