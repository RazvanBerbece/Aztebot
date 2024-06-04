package dax

type Warn struct {
	Id                int64
	UserId            string // the user by id which has been given a warn
	Reason            string // the reason for the warn
	CreationTimestamp int64  // when the warn was given to the user
}
