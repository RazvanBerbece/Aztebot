package globalConfiguration

// LOGGING
var AuditMemberJoinsInChannel = true
var AuditMemberVerificationsInChannel = true
var AuditMemberDeletesInChannel = true
var AuditRoleUpdatesInChannel = true
var AuditPromotionStateInChannel = true
var AuditPromotionMismatchesInChannel = true

// EXPERIENCE RATES
var DefaultExperienceReward_MessageSent float64 = 1.0
var DefaultExperienceReward_SlashCommandUsed float64 = 0.75
var DefaultExperienceReward_ReactionReceived float64 = 0.66
var DefaultExperienceReward_InVc float64 = 0.0075
var DefaultExperienceReward_InMusic float64 = 0.0025

var ExperienceReward_MessageSent float64 = 1.0
var ExperienceReward_SlashCommandUsed float64 = 0.75
var ExperienceReward_ReactionReceived float64 = 0.66
var ExperienceReward_InVc float64 = 0.0075
var ExperienceReward_InMusic float64 = 0.0025

// UI/UX CUSTOMISATION
var EmbedPageSize int = 10

// PROGRESSION RELATED
var SyncProgressionInMemberUpdates = true
var OrderRoleNames []string = []string{
	"🔗 Zelator",
	"📖 Theoricus",
	"📿 Philosophus",
	"🔮 Adeptus Minor",
	"〽️ Adeptus Major",
	"🧿 Adeptus Exemptus",
	"☀️ Magister Templi",
	"🧙🏼 Magus",
	"⚔️ Ipsissimus",
}
