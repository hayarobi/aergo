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
	event := &banEvent{when:time.Now(), why:why}
	bm.mutex.Lock()
	defer bm.mutex.Unlock()

	if len(addr) > 0 {
		ban, found := bm.addrMap[addr]
		if !found {
			ban = newAddrBanStatusImpl()
			bm.addrMap[addr] = ban
		}
		ban.addEvent(event)
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

func (bm *blacklistManagerImpl) IsBanned(addr string, pid peer.ID) bool {
	panic("implement me")
}

func (bm *blacklistManagerImpl) IsBannedIP(addr string) bool {
	panic("implement me")
}

func (bm *blacklistManagerImpl) IsBannedAddr(pid string) bool {
	panic("implement me")
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

