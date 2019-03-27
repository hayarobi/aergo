/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package audit

import (
	"github.com/aergoio/aergo/message"
	"github.com/aergoio/aergo/p2p/p2pcommon"
)

// TODO any logic will not be in p2pcommon package
func GetPenaltyScore(err error) p2pcommon.Penalty {
	switch penaltyEror := err.(type) {
	case message.BlamableError:
		return getPenaltyFromBlame(penaltyEror)
	default:
		return p2pcommon.PenaltyNone
	}
}

func getPenaltyFromBlame(err message.BlamableError) p2pcommon.Penalty {
	switch err.Size() {
	case message.Tiny:
		return p2pcommon.PenaltyTiny
	case message.Small:
		return p2pcommon.PenaltySmall
	case message.Normal:
		return p2pcommon.PenaltyNormal
	case message.Big:
		return p2pcommon.PenaltyBig
	case message.Severe:
		return p2pcommon.PenaltySevere
	default:
		return p2pcommon.PenaltySmall
	}
}
