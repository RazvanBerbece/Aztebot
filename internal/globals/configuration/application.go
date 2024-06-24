package globalConfiguration

// GENERAL
var VerifiedRoleName = "Aztec"
var GreetNewVerifiedUsersInChannel = false

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
var DefaultExperienceReward_InVc float64 = 0.0033
var DefaultExperienceReward_InMusic float64 = 0.00175

var ExperienceReward_MessageSent float64 = 1.0
var ExperienceReward_SlashCommandUsed float64 = 0.75
var ExperienceReward_ReactionReceived float64 = 0.66
var ExperienceReward_InVc float64 = 0.0033
var ExperienceReward_InMusic float64 = 0.00175

// COIN RATES
var DefaultCoinReward_MessageSent float64 = 2.5
var DefaultCoinReward_SlashCommandUsed float64 = 2.5
var DefaultCoinReward_ReactionReceived float64 = 7.5
var DefaultCoinReward_InVc float64 = 0.04     // 144.0 coins / hr
var DefaultCoinReward_InMusic float64 = 0.003 // 10.8 coins / hr

var CoinReward_MessageSent float64 = 2.5
var CoinReward_SlashCommandUsed float64 = 2.5
var CoinReward_ReactionReceived float64 = 7.5
var CoinReward_InVc float64 = 0.04
var CoinReward_InMusic float64 = 0.003

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
