package dax

type Timeout struct {
	Id                int64
	UserId            string // the user by id which was given a timeout
	Reason            string // the reason for the timeout
	CreationTimestamp int64  // when the timeout was given to the user
	SDuration         int    // timeout duration
}
