/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package audit

import (
	"github.com/libp2p/go-libp2p-peer"
	"sync"
	"time"
	"errors"
)

var NotFoundError = errors.New("ban status not found")
var UndefinedTime = time.Unix(0,0)

type blacklistManagerImpl struct {
    addrMap map[string]*addrBanStatusImpl
	idMap   map[peer.ID]*idBanStatusImpl

	mutex sync.Mutex
}

func NewBlacklistManager() *blacklistManagerImpl {
	return &blacklistManagerImpl{
		addrMap: make(map[string]*addrBanStatusImpl),
		idMap:   make(map[peer.ID]*idBanStatusImpl),
	}
}

func (bm *blacklistManagerImpl) AddBanScore(addr string, pid peer.ID, why string) {
	// TODO it has all same valid. make it more robust later
	now := time.Now()
	event := &banEvent{when:now, why:why}
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	if len(addr) > 0 {
		addrBan, found := bm.addrMap[addr]
		if !found {
			addrBan = newAddrBanStatusImpl()
			bm.addrMap[addr] = addrBan
		}
		addrBan.addEvent(event)
	}
	if len(pid) > 0 {
		idban, found := bm.idMap[pid]
		if !found {
			idban = newIDBanStatusImpl()
			bm.idMap[pid] = idban
		}
		idban.addEvent(event)
	}
}

func (bm *blacklistManagerImpl) IsBanned(addr string, pid peer.ID) (bool, time.Time) {
	if banned, until := bm.IsBannedAddr(addr); banned {
		return banned, until
	}
	return bm.IsBannedPeerID(pid)
}

func (bm *blacklistManagerImpl) IsBannedPeerID(peerID peer.ID) (bool, time.Time) {
	if len(peerID) > 0 {
		idban, found := bm.idMap[peerID]
		if found && time.Now().Before(idban.banUntil) {
			return true, idban.banUntil
		}
	}
	return false, UndefinedTime
}

func (bm *blacklistManagerImpl) IsBannedAddr(addr string) (bool, time.Time) {
	if len(addr) > 0 {
		addrBan, found := bm.addrMap[addr]
		if found && time.Now().Before(addrBan.banUntil) {
			return true, addrBan.banUntil
		}
	}
	return false, UndefinedTime
}

func (bm *blacklistManagerImpl) GetStatusByID(peerID peer.ID) (BanStatus, error) {
	st, found := bm.idMap[peerID]
	if !found {
		return nil, NotFoundError
	}
	return st, nil
}

func (bm *blacklistManagerImpl) GetStatusByAddr(addr string) (BanStatus, error) {
	st, found := bm.addrMap[addr]
	if !found {
		return nil, NotFoundError
	}
	return st, nil
}


func (bm *blacklistManagerImpl) NewPeerAuditor(address string, peerID peer.ID, exceedlistener ExceedListener) PeerAuditor {
	pa := NewPeerAuditor(DefaultPeerExceedThreshold, newListenWrapper(bm, exceedlistener))
	pa.peerID = peerID
	pa.ipAddress = address

	return pa
}

type listenWrapper struct {
	bm *blacklistManagerImpl
	innerListener ExceedListener
}

func newListenWrapper(bm *blacklistManagerImpl, exceedlistener ExceedListener) ExceedListener {
	return &listenWrapper{bm, exceedlistener}
}
func (lw *listenWrapper) OnExceed(auditor PeerAuditor, cause string) {
 	go lw.bm.AddBanScore(auditor.IPAddress(),auditor.PeerID(), cause)
	lw.innerListener.OnExceed(auditor, cause)
}

