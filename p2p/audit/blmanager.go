/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package audit

import (
	"errors"
	"github.com/aergoio/aergo-lib/log"
	"github.com/libp2p/go-libp2p-peer"
	"path/filepath"
	"sync"
	"time"
)

// variables that are used internally
var (
	NotFoundError = errors.New("ban status not found")
	UndefinedTime = time.Unix(0, 0)
)

const (
	blacklistFile = "blacklist.json"
	banlogFile    = "banlog.log"

	tempFileSurfix = ".tmp"

	defaultPruneInteral = time.Hour
	defaultPruneTTL = time.Hour * 24 * 730

)

type blacklistManagerImpl struct {
	logger *log.Logger

	addrMap map[string]*addrBanStatusImpl
	idMap   map[peer.ID]*idBanStatusImpl

	rwLock sync.RWMutex

	authDir string

	stopScheduler chan interface{}
}

func NewBlacklistManager(logger *log.Logger, authDir string) *blacklistManagerImpl {
	bm := &blacklistManagerImpl{
		logger:  logger,
		addrMap: make(map[string]*addrBanStatusImpl),
		idMap:   make(map[peer.ID]*idBanStatusImpl),

		authDir:       authDir,
		stopScheduler: make(chan interface{}),
	}

	return bm
}

func (bm *blacklistManagerImpl) Start() {
	bm.logger.Debug().Msg("starting up blacklist manager")
	bm.loadBlacklistFile(filepath.Join(bm.authDir, blacklistFile))
	go bm.runPruneSchedule()
}

func (bm *blacklistManagerImpl) Stop() {
	bm.logger.Debug().Msg("stopiing blacklist manager")
	bm.stopScheduler <- struct{}{}
	bm.saveBlacklistFile(filepath.Join(bm.authDir, blacklistFile))
}

func (bm *blacklistManagerImpl) AddBanScore(addr string, pid peer.ID, why string) {
	now := time.Now().Round(0)
	event := &banEvent{when: now, why: why}
	bm.rwLock.Lock()
	defer bm.rwLock.Unlock()
	bm.addAddrBanScore(addr, event)
	bm.addIDBanScore(pid, event)
}

func (bm *blacklistManagerImpl) addAddrBanScore(addr string, event *banEvent) {
	if len(addr) > 0 {
		addrBan, found := bm.addrMap[addr]
		if !found {
			addrBan = newAddrBanStatusImpl(addr)
			bm.addrMap[addr] = addrBan
		}
		addrBan.addEvent(event)
	}
}

func (bm *blacklistManagerImpl) addIDBanScore(pid peer.ID, event *banEvent) {
	if len(pid) > 0 {
		idban, found := bm.idMap[pid]
		if !found {
			idban = newIDBanStatusImpl(pid)
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
		if found && idban.Banned(time.Now()) {
			return true, idban.banUntil
		}
	}
	return false, UndefinedTime
}

func (bm *blacklistManagerImpl) IsBannedAddr(addr string) (bool, time.Time) {
	if len(addr) > 0 {
		addrBan, found := bm.addrMap[addr]
		if found && addrBan.Banned(time.Now()) {
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

func (bm *blacklistManagerImpl) runPruneSchedule() {
	pTicker := time.NewTicker(defaultPruneInteral)
	for {
		select {
		case <-pTicker.C:
			bm.pruneOldEvents()
		case <-bm.stopScheduler:
			break
		}
	}
}

func (bm *blacklistManagerImpl) pruneOldEvents() {
	bm.rwLock.Lock()
	defer bm.rwLock.Unlock()
	// pruning is not applied to banned peer
	pruneDelay := time.Now().Add( defaultPruneTTL * -1)
	for _, bs := range bm.addrMap {
		if !bs.Banned(pruneDelay) {
			bs.PruneOldEvents(pruneDelay)
		}
	}
	for _, bs := range bm.idMap {
		if !bs.Banned(pruneDelay) {
			bs.PruneOldEvents(pruneDelay)
		}
	}
}

type listenWrapper struct {
	bm            *blacklistManagerImpl
	innerListener ExceedListener
}

func newListenWrapper(bm *blacklistManagerImpl, exceedlistener ExceedListener) ExceedListener {
	return &listenWrapper{bm, exceedlistener}
}

func (lw *listenWrapper) OnExceed(auditor PeerAuditor, cause string) {
	go lw.bm.AddBanScore(auditor.IPAddress(), auditor.PeerID(), cause)
	lw.innerListener.OnExceed(auditor, cause)
}
