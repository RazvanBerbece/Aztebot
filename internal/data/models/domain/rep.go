package domain

import "time"

type GivenRep struct {
	To        string
	Timestamp time.Time
}

func IdInAuthorRepTargetList(targetId string, repList []GivenRep) bool {
	for _, rep := range repList {
		if rep.To == targetId {
			return true
		}
	}
	return false
}
