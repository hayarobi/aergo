/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package audit

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
	category PenaltyCategory
	score    PenaltyPoint
}

func (p Penalty) String() string {
	return fmt.Sprintf("%v/%f",p.category,p.score)
}
