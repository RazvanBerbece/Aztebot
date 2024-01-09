package dataModels

type Warn struct {
	UserId    string // the user by id which has been given a warn
	Reason    string // the reason for the warn
	Timestamp int64  // when the warn was given to the user
}
