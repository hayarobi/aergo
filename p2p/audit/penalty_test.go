/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package audit

import (
	"github.com/aergoio/aergo/message"
	"github.com/aergoio/aergo/types"
	"reflect"
	"testing"
)


func TestGetPenaltyScore(t *testing.T) {

	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want Penalty
	}{
		{"TBlameSevere", args{message.NewBlamableError(message.Severe,"")}, PenaltySevere},
		{"TBlameBig", args{message.NewBlamableError(message.Big,"")}, PenaltyBig},
		{"TBlameNormal", args{message.NewBlamableError(message.Normal,"")}, PenaltyNormal},
		{"TBlameTiny", args{message.NewBlamableError(message.Tiny,"")}, PenaltyTiny},
		{"TNotBlame", args{types.ErrSignNotMatch}, PenaltyNone},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPenaltyScore(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPenaltyScore() = %v, want %v", got, tt.want)
			}
		})
	}
}
