/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package audit

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestDefaultAuditor_AddScore(t *testing.T) {
	//t.Skip("Too long time test")
	type fields struct {
		threshold      float64
	}
	type args struct {
		category PenaltyCategory
		score    float64
		count    int
		interval time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args

		wantListenerCalled bool
	}{
		{"TSmallPerm",fields{10000}, args{Permanent,1000,10, time.Millisecond}, false},
		{"TExPerm",fields{10000}, args{Permanent,1000,11, time.Millisecond}, true},
		//{"TSmallLong",fields{10000}, args{LongTerm,910,11, time.Second>>2}, false},
		//{"TExLong",fields{10000}, args{LongTerm,911,11, time.Second>>2}, true},
		//{"TSmallShort",fields{10000}, args{ShortTerm,980,11, time.Second>>2}, false},
		//{"TExShort",fields{10000}, args{ShortTerm,1000,11, time.Second>>2}, true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listener := &TestExceedListener{}
			a := NewPeerAuditor(tt.fields.threshold, listener)

			for i:=0; i<tt.args.count; i++ {
				a.AddScore(tt.args.category, tt.args.score)
				time.Sleep(tt.args.interval)
			}

			if tt.wantListenerCalled != listener.Called() {
				t.Errorf("exceed listener OnExceed %v, want %v",listener.Called(),tt.wantListenerCalled)
			}
		})
	}
}

type TestExceedListener struct {
	called int32
}

func (l *TestExceedListener) OnExceed(auditor PeerAuditor, cause string) {
	atomic.StoreInt32(&l.called,1)
}

func (l *TestExceedListener) Called() bool {
	return atomic.LoadInt32(&l.called) == 1
}

