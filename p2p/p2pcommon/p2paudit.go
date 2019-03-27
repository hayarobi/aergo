/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package p2pcommon

import "fmt"

type PenaltyCategory int32
type PenaltyPoint float64

const (
	ShortTerm PenaltyCategory = iota
	LongTerm
	Permanent
)

//go:generate stringer -type=PenaltyCategory

type Penalty struct {
	Category PenaltyCategory
	Score    PenaltyPoint
}

func (p Penalty) String() string {
	return fmt.Sprintf("%v/%f",p.Category,p.Score)
}


// Pre-defined penalties
var (
	PenaltyNone   = Penalty{Permanent, 0.0}
	PenaltySevere = Penalty{Permanent, 100001.0}
	PenaltyBig    = Penalty{Permanent, 50000.0}
	PenaltyNormal = Penalty{LongTerm, 10000.0}
	PenaltySmall  = Penalty{ShortTerm, 1000.0}
	PenaltyTiny   = Penalty{ShortTerm, 100.0}
)

