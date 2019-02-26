/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package audit

import (
	"strconv"
	"testing"
	"time"
)

func Test_banStatusImpl_addEvent(t *testing.T) {
	// added events must be sorted by valid time.
	type args struct {
		kickTimes int
	}
	tests := []struct {
		name string
		args args
		wantBan bool
		wantDuration time.Duration
	}{
		{"TBeforeBan", args{2}, false, 0},
		{"TFirstBan", args{3}, true, BanDurations[2]},
		{"TooManyBan", args{15}, true, BanDurations[len(BanDurations)-1]},
		{"TooManyBan2", args{20}, true, BanDurations[len(BanDurations)-1]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseTime := time.Now()
			bs := newIDBanStatusImpl()
			for i:=0 ; i< tt.args.kickTimes ; i++{
				ev := &banEvent{why:strconv.Itoa(i),when:time.Now()}
				bs.addEvent(ev)
			}
			gotDuration := bs.banUntil.Sub(baseTime)
			if tt.wantBan != (gotDuration > 0) {
				t.Errorf("banStatusImpl.addEvent() ban %v , want %v ",gotDuration>0, tt.wantBan)
			}
			if tt.wantBan {
				diff := gotDuration - tt.wantDuration
				if diff < 0 {
					diff *= -1
				}
				if diff > time.Second*4 {
					t.Errorf("banStatusImpl.addEvent() banDuration %v , want %v ",gotDuration, tt.wantDuration)
				}
			}
		})
	}
}

func Test_newIDBanStatusImpl(t *testing.T) {
	bs := newIDBanStatusImpl()
	if bs.ValidUntil() != UndefinedTime {
		t.Errorf("wrong initialization in validuntil")
	}
	if len(bs.Events()) > 0 {
		t.Errorf("wrong initialization in event count")
	}
}

func Test_newAddrBanStatusImpl(t *testing.T) {
	bs := newAddrBanStatusImpl()
	if bs.ValidUntil() != UndefinedTime {
		t.Errorf("wrong initialization in validuntil")
	}
	if len(bs.Events()) > 0 {
		t.Errorf("wrong initialization in event count")
	}
}
