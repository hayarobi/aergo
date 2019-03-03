/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package audit

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-peer"
)

var dummyID, _ = peer.IDB58Decode("16Uiu2HAmFqptXPfcdaCdwipB2fhHATgKGVFVPehDAPZsDKSU7jRm")

func Test_banStatusImpl_addEvent(t *testing.T) {
	// added events must be sorted by valid time.
	type args struct {
		kickTimes int
	}
	tests := []struct {
		name         string
		args         args
		wantBan      bool
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
			bs := newIDBanStatusImpl(dummyID)
			for i := 0; i < tt.args.kickTimes; i++ {
				ev := &banEvent{why: strconv.Itoa(i), when: time.Now()}
				bs.addEvent(ev)
			}
			gotDuration := bs.banUntil.Sub(baseTime)
			if tt.wantBan != (gotDuration > 0) {
				t.Errorf("banStatusImpl.addEvent() ban %v , want %v ", gotDuration > 0, tt.wantBan)
			}
			if tt.wantBan {
				diff := gotDuration - tt.wantDuration
				if diff < 0 {
					diff *= -1
				}
				if diff > time.Second*4 {
					t.Errorf("banStatusImpl.addEvent() banDuration %v , want %v ", gotDuration, tt.wantDuration)
				}
			}
		})
	}
}

func Test_newIDBanStatusImpl(t *testing.T) {
	bs := newIDBanStatusImpl(dummyID)
	if bs.BanUntil() != UndefinedTime {
		t.Errorf("wrong initialization in validuntil")
	}
	if len(bs.Events()) > 0 {
		t.Errorf("wrong initialization in event count")
	}
}

func Test_newAddrBanStatusImpl(t *testing.T) {
	bs := newAddrBanStatusImpl(dummyAddress)
	if bs.BanUntil() != UndefinedTime {
		t.Errorf("wrong initialization in validuntil")
	}
	if len(bs.Events()) > 0 {
		t.Errorf("wrong initialization in event count")
	}
}

func Test_newStatusImpl(t *testing.T) {
	var src []int = nil
	abc := make([]int, 0)
	copy(abc, src)
	fmt.Println("Result ", abc)

	type args struct {
		initialScore int
		banTime      time.Time
		prevEvents   []BanEvent
	}
	tests := []struct {
		name string
		args args
		wantScore int
		wantEvCnt int
	}{
		{"TNil", args{0, time.Time{}, nil}, 0, 0},
		{"TPreScore", args{31, time.Time{}, nil}, 31, 0},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newStatusImpl(tt.args.initialScore, tt.args.banTime, tt.args.prevEvents);
			if got.banScore != tt.wantScore {
				t.Errorf("newStatusImpl() score = %v, want %v", got.banScore, tt.wantScore)
			}
			if len(got.events) != tt.wantEvCnt {
				t.Errorf("newStatusImpl() score = %v, want %v", len(got.events), tt.wantEvCnt)
			}
		})
	}

}


func Test_banStatusImpl_PruneOldEvents(t *testing.T) {
	preCnt := 20
	oldestTime := time.Now().Round(0)

	type args struct {
		pruneTime time.Time
	}
	tests := []struct {
		name   string
		args   args
		want   int

	}{
		{"TAllNew", args{oldestTime.Add(-1 * time.Second)},0},
		{"TFirstRemoved", args{oldestTime.Add(time.Second)},1},
		{"TLastOne", args{oldestTime.Add(time.Minute*time.Duration(preCnt-1)-time.Second)},19},
		{"TAllOld", args{oldestTime.Add(time.Minute*time.Duration(preCnt-1))},20},
		{"TFarOld", args{oldestTime.Add(time.Hour*1000)},20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := newStatusImpl(0, UndefinedTime, nil)
			for i:=0; i<preCnt; i++ {
				bs.addEvent(&banEvent{why: strconv.Itoa(i), when: oldestTime.Add(time.Minute*time.Duration(i))})
			}

			if got := bs.PruneOldEvents(tt.args.pruneTime); got != tt.want {
				t.Errorf("banStatusImpl.PruneOldEvents() = %v, want %v", got, tt.want)
			}
			if len(bs.events) != (preCnt-tt.want) {
				t.Errorf("banStatusImpl.PruneOldEvents() remained events = %v, want %v", len(bs.events), preCnt-tt.want)
			}
		})
	}
}
