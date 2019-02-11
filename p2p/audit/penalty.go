/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package audit

import (
	"github.com/aergoio/aergo/message"
)

var (
	PenaltyNone   = Penalty{Permanent, 0.0}
	PenaltySevere = Penalty{Permanent, 100001.0}
	PenaltyBig    = Penalty{Permanent, 50000.0}
	PenaltyNormal = Penalty{LongTerm , 10000.0}
	PenaltySmall  = Penalty{ShortTerm, 1000.0}
	PenaltyTiny   = Penalty{ShortTerm, 100.0}
)

func GetPenaltyScore(err error) Penalty {
	switch penaltyEror := err.(type) {
	case message.BlamableError:
		return getPenaltyFromBlame(penaltyEror)
	default:
		return PenaltyNone
	}
}

func getPenaltyFromBlame(err message.BlamableError) Penalty {
	switch err.Size() {
	case message.Tiny:
		return PenaltyTiny
	case message.Small:
		return PenaltySmall
	case message.Normal:
		return PenaltyNormal
	case message.Big:
		return PenaltyBig
	case message.Severe:
		return PenaltySevere
	default:
		return PenaltySmall
	}
}
