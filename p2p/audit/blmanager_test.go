/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package audit

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-peer"
)

func TestNewBlacklistManager(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"T1"},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewBlacklistManager(nil, "")
			if got.addrMap == nil {
				t.Errorf("NewBlacklistManager() fields not initialized %v","addrMap")
			}
			if got.idMap == nil {
				t.Errorf("NewBlacklistManager() fields not initialized %v","addrMap")
			}
		})
	}
}

func Test_blacklistManagerImpl_AddBanScore(t *testing.T) {
	addr1 := "123.45.67.89"
	id1, _ := peer.IDB58Decode("16Uiu2HAmFqptXPfcdaCdwipB2fhHATgKGVFVPehDAPZsDKSU7jRm")
	addrother := "8.8.8.8"
	idother, _ := peer.IDB58Decode("16Uiu2HAmU8Wc925gZ5QokM4sGDKjysdPwRCQFoYobvoVnyutccCD")
	// this test runs incrementally; it keep state of prev testcases.
	type args struct {
		addr string
		pid  peer.ID
		why  string
	}
	tests := []struct {
		name   string
		args   args
		wantAddrCnt int
		wantIdCount int
	}{
		{"TOne",args{addr1, id1, "first"},1,1},
		{"TDiffOne",args{addrother, idother, "diff"},1,1},
		{"TSameAddr",args{addr1, idother, "addr"},2,1},
		{"TSameId",args{addrother, id1, "id"},2,2},
		{"TBoth",args{addr1, id1, "both"},3,3},
	}
	bm := NewBlacklistManager(nil, "")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm.AddBanScore(tt.args.addr, tt.args.pid, tt.args.why)
			if _, found := bm.addrMap[addr1] ; !found {
				t.Errorf("AddBanScore() not added addr ")
			}
			if _, found := bm.idMap[id1] ; !found {
				t.Errorf("AddBanScore() not added id ", )
			}
			addrStat := bm.addrMap[addr1]
			if len(addrStat.events) != tt.wantAddrCnt {
				t.Errorf("blacklistManagerImpl.AddBanScore() addr cnt = %v, want %v",len(addrStat.events), tt.wantAddrCnt)
			}
			idStat := bm.idMap[id1]
			if len(idStat.events) != tt.wantIdCount {
				t.Errorf("blacklistManagerImpl.AddBanScore() id cnt = %v, want %v",len(idStat.events), tt.wantAddrCnt)
			}
		})
	}
}

func Test_blacklistManagerImpl_IsBanned(t *testing.T) {
	addr1 := "123.45.67.89"
	id1, _ := peer.IDB58Decode("16Uiu2HAmFqptXPfcdaCdwipB2fhHATgKGVFVPehDAPZsDKSU7jRm")
	addrother := "8.8.8.8"
	idother, _ := peer.IDB58Decode("16Uiu2HAmU8Wc925gZ5QokM4sGDKjysdPwRCQFoYobvoVnyutccCD")

	type args struct {
		addr string
		pid  peer.ID
	}
	tests := []struct {
		name   string
		args   args
		want   bool
	}{
		{"TNotFound",args{addrother, idother},false},
		{"TFoundAddr",args{addr1, idother},true},
		{"TFoundId",args{addrother, id1},true},
		{"TFoundBoth",args{addr1, id1},true},
	}
	b := NewBlacklistManager(nil, "")
	for i:=0 ; i < 10 ; i++ {
		b.AddBanScore(addr1, id1, "test "+strconv.Itoa(i))
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := b.IsBanned(tt.args.addr, tt.args.pid); got != tt.want {
				t.Errorf("blacklistManagerImpl.IsBanned() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_blacklistManagerImpl_IsBannedAddr(t *testing.T) {
	type fields struct {
		addrMap map[string]*addrBanStatusImpl
		ipMap   map[peer.ID]*idBanStatusImpl
	}
	type args struct {
		pid string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &blacklistManagerImpl{
				addrMap: tt.fields.addrMap,
				idMap:   tt.fields.ipMap,
			}
			if got, _ := b.IsBannedAddr(tt.args.pid); got != tt.want {
				t.Errorf("blacklistManagerImpl.IsBannedAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_blacklistManagerImpl_GetStatus(t *testing.T) {
	type fields struct {
		addrMap map[string]*addrBanStatusImpl
		ipMap   map[peer.ID]*idBanStatusImpl
	}
	type args struct {
		addr string
		pid  peer.ID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    BanStatus
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &blacklistManagerImpl{
				addrMap: tt.fields.addrMap,
				idMap:   tt.fields.ipMap,
			}
			got, err := b.GetStatusByAddr(tt.args.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("blacklistManagerImpl.GetStatusByAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("blacklistManagerImpl.GetStatusByAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_blacklistManagerImpl_NewPeerAuditor(t *testing.T) {
	// check auditor is made, and correctly observe peer's exceed notice

	addr1 := "123.45.67.89"
	id1, _ := peer.IDB58Decode("16Uiu2HAmFqptXPfcdaCdwipB2fhHATgKGVFVPehDAPZsDKSU7jRm")
	type args struct {
	}
	tests := []struct {
		name   string
		args   args
		want   PeerAuditor
	}{
		{"T1",args{}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dl := &dummyListener{}

			bm := NewBlacklistManager(nil, "")
			got := bm.NewPeerAuditor(addr1, id1, dl)
			if  got.PeerID() != id1  {
				t.Errorf("blacklistManagerImpl.NewPeerAuditor() = %v, want %v",  got.PeerID() , id1)
			}
			if  got.IPAddress() != addr1 {
				t.Errorf("blacklistManagerImpl.NewPeerAuditor() = %v, want %v",  got.IPAddress() , addr1)
			}

			got.AddPenalty(PenaltySevere)
			time.Sleep(time.Millisecond<<4)

			if !dl.called {
				t.Errorf("exceed event is expected, but was not.", )
			}
			stat, _ := bm.GetStatusByID(id1)
			if len(stat.Events()) == 0 {
				t.Errorf("ban stat was not increased for id %v", id1 )
			}
			stataddr, _ := bm.GetStatusByAddr(addr1)
			if len(stataddr.Events()) == 0 {
				t.Errorf("ban stat was not increased for addr %v", addr1 )
			}
		})
	}
}

type dummyListener struct {
	called bool
}

func (l *dummyListener) OnExceed(auditor PeerAuditor, cause string) {
	l.called = true
	fmt.Printf("got exceed notice: %v",cause)
}
