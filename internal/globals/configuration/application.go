package globalConfiguration

// GENERAL
const DefaultVerifiedRoleId = 1
const DefaultVerifiedRoleName = "Aztec"
const GreetNewVerifiedUsersInChannel = false

const ActivityId_MessageSend = "msg_send"
const ActivityId_SlashCommandUse = "slash_command_used"
const ActivityId_ReactionReceived = "react_recv"
const ActivityId_TimeInVc = "in_vc"
const ActivityId_TimeInMusic = "in_music"

// LOGGING
const AuditMemberJoinsInChannel = true
const AuditMemberVerificationsInChannel = true
const AuditMemberDeletesInChannel = true
const AuditRoleUpdatesInChannel = true
const AuditPromotionStateInChannel = true
const AuditPromotionMismatchesInChannel = true

// EXPERIENCE RATES
// defaults (do not change unless for finetuning)
const DefaultExperienceReward_MessageSent float64 = 1.0
const DefaultExperienceReward_SlashCommandUsed float64 = 0.75
const DefaultExperienceReward_ReactionReceived float64 = 0.66
const DefaultExperienceReward_InVc float64 = 0.0033
const DefaultExperienceReward_InMusic float64 = 0.00175

// note, below are actively used in reward processing
var ExperienceReward_MessageSent float64 = 1.0
var ExperienceReward_SlashCommandUsed float64 = 0.75
var ExperienceReward_ReactionReceived float64 = 0.66
var ExperienceReward_InVc float64 = 0.0033
var ExperienceReward_InMusic float64 = 0.00175

// COIN RATES
// defaults (do not change unless for finetuning)
const DefaultCoinReward_MessageSent float64 = 2.5
const DefaultCoinReward_SlashCommandUsed float64 = 2.5
const DefaultCoinReward_ReactionReceived float64 = 7.5
const DefaultCoinReward_InVc float64 = 0.04     // 144.0 coins / hr
const DefaultCoinReward_InMusic float64 = 0.003 // 10.8 coins / hr

// note, below are actively used in reward processing
var CoinReward_MessageSent float64 = 2.5
var CoinReward_SlashCommandUsed float64 = 2.5
var CoinReward_ReactionReceived float64 = 7.5
var CoinReward_InVc float64 = 0.04
var CoinReward_InMusic float64 = 0.003

// UI/UX CUSTOMISATION
const EmbedPageSize int = 10

// PROGRESSION RELATED
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
