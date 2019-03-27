/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package audit

import (
	"github.com/aergoio/aergo/p2p/p2pcommon"
	"github.com/libp2p/go-libp2p-peer"
	"time"
)

type dummyBlacklistManager struct {

}

func newDummyBlacklistManager() BlacklistManager {
	return &dummyBlacklistManager{}
}


func (*dummyBlacklistManager) Start() {
}

func (*dummyBlacklistManager) Stop() {
}

func (*dummyBlacklistManager) NewPeerAuditor(address string, peerID peer.ID, exceedlistener ExceedListener) PeerAuditor {
	return &dummyAuditor{peerID:peerID, ipAddress:address}
}

func (*dummyBlacklistManager) IsBanned(addr string, pid peer.ID) (bool, time.Time) {
	return false, UndefinedTime
}

func (*dummyBlacklistManager) IsBannedPeerID(peerID peer.ID) (bool, time.Time) {
	return false, UndefinedTime
}

func (*dummyBlacklistManager) IsBannedAddr(addr string) (bool, time.Time) {
	return false, UndefinedTime
}

func (*dummyBlacklistManager) GetStatusByID(peerID peer.ID) (BanStatus, error) {
	return nil, NotFoundError
}

func (*dummyBlacklistManager) GetStatusByAddr(addr string) (BanStatus, error) {
	return nil, NotFoundError
}

func (*dummyBlacklistManager) Summary() (map[string]interface{}) {
	sum := make(map[string]interface{})
	idBan := make(map[string] interface{})
	addrBan := make(map[string] interface{})

	sum["bannedID"] = idBan
	sum["bannedAddr"] = addrBan
	return sum
}

type dummyAuditor struct {
	peerID peer.ID
	ipAddress string
}

func (da *dummyAuditor) PeerID() peer.ID {
	return da.peerID
}

func (da *dummyAuditor) IPAddress() string {
	return da.ipAddress
}

func (*dummyAuditor) AddPenalty(penalty p2pcommon.Penalty) bool {
	return true
}

func (*dummyAuditor) AddScore(category p2pcommon.PenaltyCategory, score float64) bool {
	return true
}

func (*dummyAuditor) Threshold() float64 {
	return DefaultPeerExceedThreshold
}

func (*dummyAuditor) CurrentScore(category p2pcommon.PenaltyCategory) float64 {
	return 0.0
}

func (*dummyAuditor) ScoreSum() float64 {
	return 0.0
}
