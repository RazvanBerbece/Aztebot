package dax

type ArchivedTimeout struct {
	Id              int64
	UserId          string // the user by id which was given a timeout
	Reason          string // the reason for the timeout
	ExpiryTimestamp int64  // when the timeout expired
}
