package audit

import (
	"encoding/csv"
	"errors"
	"github.com/aergoio/aergo/cmd/aergocli/util/encoding/json"
	"github.com/libp2p/go-libp2p-peer"
	"io"
	"net"
	"os"
	"path/filepath"
	"time"
)

const LogFieldsCount = 6
const timeFormat = "2019-01-01T00:00:00.001Z"

var (
	invalidBanLogFormat = errors.New("invalid band log format")
	dummyAddress = "192.168.1.3"
)

func ReadCSVLog(filepath string) (map[string]*addrBanStatusImpl, map[peer.ID]*idBanStatusImpl, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, nil, err
	}

	addrMap := make(map[string]*addrBanStatusImpl)
	idMap := make(map[peer.ID]*idBanStatusImpl)
	r := csv.NewReader(file)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		if len(record) != LogFieldsCount {
			return nil, nil, invalidBanLogFormat
		}
		when, err := time.Parse(timeFormat, record[0])
		if err != nil {
			return nil, nil, invalidBanLogFormat
		}
		why := record[1]
		addr := record[2]
		ip := net.ParseIP(addr)
		if err != nil {
			return nil, nil, err
		}
		if ip != nil {
			addrBan := addrMap[addr]
			if addrBan == nil {
				addrBan := newAddrBanStatusImpl(dummyAddress)
				addrMap[addr] = addrBan
			}
			addrBan.addEvent(&banEvent{when: when, why: why})
		}

		//id := record[4]

	}

	return addrMap, idMap, nil
}

func (bm *blacklistManagerImpl) loadBlacklistFile(filePath string) {
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		bm.logger.Warn().Err(err).Msg("Failed to load blacklist file")
		return
	}
	bm.rwLock.RLock()
	defer bm.rwLock.RUnlock()

	jr := json.NewDecoder(file)
	mc := &blacklistStats{}
	err = jr.Decode(mc)
	if err != nil {
		bm.logger.Warn().Err(err).Msg("Failed to decode blacklist file")
		return
	}

	bm.addrMap = make(map[string]*addrBanStatusImpl, len(bm.addrMap))
	bm.idMap = make(map[peer.ID]*idBanStatusImpl, len(bm.idMap))
	for _, b := range mc.Addrs {
		status := &addrBanStatusImpl{addr: b.Addr, banStatusImpl:banStatusImpl{banScore:b.BanScore, banUntil: b.BanUntil}}
		status.events = make([]BanEvent, len(b.Events))
		for j, e := range b.Events {
			status.events[j] = &banEvent{when:e.When, why:e.Why}
		}
		bm.addrMap[b.Addr] = status
	}
	for _, b := range mc.Ids {
		id, err := peer.IDB58Decode(b.Id)
		if err != nil {
			bm.logger.Warn().Str("id", b.Id).Err(err).Msg("failed to decode id")
			continue
		}
		status := &idBanStatusImpl{id:id, banStatusImpl:banStatusImpl{banScore:b.BanScore, banUntil: b.BanUntil}}
		status.events = make([]BanEvent, len(b.Events))
		for j, e := range b.Events {
			status.events[j] = &banEvent{when:e.When, why:e.Why}
		}
		bm.idMap[id] = status
	}

	bm.appendWAL()
}

func (bm *blacklistManagerImpl) appendWAL() {
	// TODO read WAL and apply to map
}

func (bm *blacklistManagerImpl) saveBlacklistFile(filePath string) {
	parent := filepath.Dir(filePath)
	if err := os.MkdirAll(parent, 0755) ; err!=nil {
		bm.logger.Info().Err(err).Msg("Failed to create auth directory")
	}

	tmpFile := filePath+tempFileSurfix
	file, err := os.OpenFile(tmpFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		bm.logger.Warn().Err(err).Msg("Failed to save temporary file for blacklist")
	}
	bm.rwLock.Lock()
	defer bm.rwLock.Unlock()
	jr := json.NewEncoder(file)
	mc := convertStat(time.Now(), bm)
	err = jr.Encode(&mc)
	if err != nil {
		panic("err when encode:" + err.Error())
	}
	if err = os.Rename(tmpFile, filePath); err != nil {
		bm.logger.Warn().Err(err).Msg("Failed to save blacklist file. can't replace old file")
	}
}

func convertStat(savetime time.Time, bm *blacklistManagerImpl) *blacklistStats {
	addrs := make([]addrBanInfo, len(bm.addrMap))
	i:=0
	for k, b := range bm.addrMap {
		addrs[i] = addrBanInfo{Addr:k, BanScore: b.banScore, BanUntil:b.banUntil}
		addrs[i].Events = make([]eventInfo,len(b.events))
		for j, e := range b.events {
			addrs[i].Events[j] = eventInfo{When:e.When(), Why:e.Why()}
		}
		i++
	}

	i = 0
	ids := make([]idBanInfo, len(bm.idMap))
	for k, b := range bm.idMap {
		ids[i] = idBanInfo{Id: k.Pretty(), BanScore: b.banScore, BanUntil:b.banUntil}
		ids[i].Events = make([]eventInfo,len(b.events))
		for j, e := range b.events {
			ids[i].Events[j] = eventInfo{When:e.When(), Why:e.Why()}
		}
		i++
	}

	return &blacklistStats{SaveTime:savetime, Addrs:addrs, Ids:ids}
}

type blacklistStats struct {
	SaveTime time.Time
	Addrs []addrBanInfo
	Ids []idBanInfo
}

type addrBanInfo struct {
	Addr string
	BanScore int
	BanUntil time.Time
	Events []eventInfo
}


type idBanInfo struct {
	Id string
	BanScore int
	BanUntil time.Time
	Events []eventInfo
}

type eventInfo struct {
	When time.Time
	Why string
}

