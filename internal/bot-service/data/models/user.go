package dataModels

type User struct {
	Id                int
	DiscordTag        string
	UserId            string
	CurrentRoleIds    string
	CurrentCircle     string
	CurrentInnerOrder *int
	CurrentLevel      int
	CurrentExperience float64
	CreatedAt         *int64
	Gender            int8
}
