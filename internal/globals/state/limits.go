package globalState

import "time"

var LastUsedTopTimestamp = time.Now().Add(-60 * time.Minute)
var LastUsedTop5sTimestamp = time.Now().Add(-60 * time.Minute)
var LastUsedMonthlyLeaderboardTimestamp = time.Now().Add(-60 * time.Minute)

var DynamicChannelsCount int = 0

var MaxTotalDynamicChannelsCount int = 10
var MaxDynamicChannelPerCategoryCount int = 5
