package audit

import (
	"bytes"
	"fmt"
	"github.com/aergoio/aergo/cmd/aergocli/util/encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-peer"
)

func Test_blacklistManagerImpl_saveBlacklistFile(t *testing.T) {
	const (
		addr1 = "192.168.0.12"
		addr2 = "172.21.0.3"
		addr3 = "1.2.3.4"
	)
	var(
		id1, _ = peer.IDB58Decode("16Uiu2HAmFqptXPfcdaCdwipB2fhHATgKGVFVPehDAPZsDKSU7jRm")
		id2, _ = peer.IDB58Decode("16Uiu2HAmU8Wc925gZ5QokM4sGDKjysdPwRCQFoYobvoVnyutccCD")
	)
	tempFile := "test1.json"
	if _, err := os.Stat(tempFile); err == nil {
		t.Skip("file path for test already exists. ", tempFile)
	}
	defer func() {
		if t.Failed() {
			b, err := ioutil.ReadFile(tempFile)
			if err != nil {
				fmt.Println("failed to dump temp file ",err)
			} else {
				fmt.Println("Saved json: ")
				var out bytes.Buffer
				json.Indent(&out, b, "","  ")
				fmt.Println(out.String())
			}
		}
		err := os.Remove(tempFile);	if err != nil {	t.Log("Failed to remote temp file ", err.Error())} }()

	// make sameple bm
	type jj struct {
		addr string
		id peer.ID
		ev *banEvent
	}
	ins := []jj{
		// events time is wall clock reading
		{addr1,id1, &banEvent{time.Now().Add(- time.Hour * 24).Round(0),"001"}},
		{addr1,id1, &banEvent{time.Now().Add(- time.Hour * 12).Round(0),"002"}},
		{addr2,id1, &banEvent{time.Now().Add(- time.Hour * 4).Round(0),"003"}},
		{addr2,id2, &banEvent{time.Now().Add(- time.Hour * 2).Round(0),"004"}},
		{addr3,id1, &banEvent{time.Now().Add(- time.Hour * 1).Round(0),"005"}},
	}
	bm := NewBlacklistManager(nil, "")
	for _, in := range ins {
		bm.addAddrBanScore(in.addr, in.ev)
		bm.addIDBanScore(in.id, in.ev)
	}

	type args struct {
		filePath string
	}
	tests := []struct {
		name   string
		args   args
	}{
		{"T1", args{tempFile} },
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bm.saveBlacklistFile(tt.args.filePath)

			other := NewBlacklistManager(nil, "")

			other.loadBlacklistFile(tt.args.filePath)

			if len(bm.addrMap) != len(other.addrMap) {
				t.Errorf("Save or Load failed: addrMap length is differ %v , want %v ",other.addrMap,bm.addrMap)
			} else {
				for k, s := range bm.addrMap {
					oStat, found := other.addrMap[k]
					if !found {
						t.Errorf("Save or Load failed: addr %v is missing  ",k)
						break
					}
					if !isEqualAddrBan(s, oStat) {
						t.Errorf("Save or Load failed: got %v , want %v",oStat, s)
						break
					}
				}
			}

			if len(bm.idMap) != len(other.idMap) {
				t.Errorf("Save or Load failed: idMap length is differ %v , want %v ",other.idMap,bm.idMap)
			} else {
				for k, s := range bm.idMap {
					oStat, found := other.idMap[k]
					if !found {
						t.Errorf("Save or Load failed: addr %v is missing  ",k)
						break
					}
					if !isEqualIdBan(s, oStat) {
						t.Errorf("Save or Load failed: got %v , want %v",oStat, s)
						break
					}
				}
			}
		})
	}
}

func isEqualAddrBan(a, b *addrBanStatusImpl) bool {
	if a.addr != b.addr {
		return false
	}
	return isEqualBanStatus(&a.banStatusImpl,&b.banStatusImpl)
}

func isEqualIdBan(a, b *idBanStatusImpl) bool {
	if a.id != b.id {
		return false
	}
	return isEqualBanStatus(&a.banStatusImpl,&b.banStatusImpl)
}
func isEqualBanStatus(a, b *banStatusImpl) bool {
	if a.banUntil != b.banUntil {
		return false
	}
	if a.banScore != b.banScore {
		return false
	}

	if len(a.events) != len(b.events) {
		return false
	}
	EV:
	for _, ae := range a.events {
		for _, be := range b.events {
			if isEqualEvent(ae, be) {
				continue EV
			}
		}
		return false
	}
	return true
}
func isEqualEvent(a, b BanEvent) bool {
	return a.When() == b.When() && a.Why() == b.Why()
}