/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package audit

import (
	"fmt"
	"github.com/libp2p/go-libp2p-peer"
	"time"
)

// banStatusImpl contains common properties and methods of implementations of BanStatus interface
type banStatusImpl struct {
	// events is history of events which affect ban scores
	events []BanEvent

	// banScore is cumulated value of not-expired ban events
	banScore int

	banUntil time.Time
}

func newStatusImpl(initialScore int, banUntil time.Time, prevEvents []BanEvent) banStatusImpl {
	tmp := make([]BanEvent, len(prevEvents))
	if len(tmp) > 0 {
		copy(tmp, prevEvents)
	}

	return banStatusImpl{banScore:initialScore, events:tmp, banUntil:banUntil}
}

func (bs *banStatusImpl) BanUntil() time.Time {
	return bs.banUntil
}

func (bs *banStatusImpl) Banned(refTime time.Time) bool {
	return refTime.Before(bs.banUntil)
}

func (bs *banStatusImpl) Events() []BanEvent {
	return bs.events
}

// PruneOldEvents remove expired ban events and re-caculate ban score.
func (bs *banStatusImpl) PruneOldEvents(pruneTime time.Time) int {
	events := len(bs.events)
	if events == 0 {
		return 0
	}
	// bs.events are sorted by time
	idx := 0
	for ;idx < events; idx++ {
		if bs.events[idx].When().After(pruneTime) {
			break
		}
	}
	if idx > 0 {
		bs.events = bs.events[idx:]
	}
	return idx
}


func (bs *banStatusImpl)addEvent(ev BanEvent) {
	// this method should be inside lock
	// recalculate total score
	// add last event
	bs.events = append(bs.events, ev)
	bs.banScore++

	// set clock to monotonic clocks since there is some trouble when marshal/unmarshal time with wall clock.
	bs.updateStats(time.Now().Round(0))
}

func (bs *banStatusImpl) updateStats(now time.Time) {
	durationIdx := bs.banScore - 1
	if durationIdx >= len(BanDurations) {
		durationIdx = len(BanDurations)-1
	}
	if BanDurations[durationIdx] == 0 {
		bs.banUntil = UndefinedTime
	} else {
		bs.banUntil = now.Add(BanDurations[durationIdx])
	}
}

type idBanStatusImpl struct {
	banStatusImpl
	id peer.ID
}

func newIDBanStatusImpl(id peer.ID) *idBanStatusImpl {
	return &idBanStatusImpl{banStatusImpl: newStatusImpl(0, UndefinedTime,nil),id:id}
}


func (bs *idBanStatusImpl) ID() string {
	return bs.id.Pretty()
}

func (bs *idBanStatusImpl) String() string {
	return fmt.Sprintf("id:%s , score:%d, until:%v , evs: %v ",bs.id.Pretty(),bs.banScore,bs.banUntil, bs.events)
}

type addrBanStatusImpl struct {
	banStatusImpl
	addr string
}

func newAddrBanStatusImpl(addr string) *addrBanStatusImpl {
	return &addrBanStatusImpl{banStatusImpl:newStatusImpl(0, UndefinedTime,nil), addr:addr}
}

func (bs *addrBanStatusImpl) ID() string {
	return bs.addr
}

func (bs *addrBanStatusImpl) String() string {
	return fmt.Sprintf("id:%s , score:%d, until:%v , evs: %v ",bs.addr,bs.banScore,bs.banUntil, bs.events)
}


type banEvent struct {
	when time.Time
	why string
}

func (e *banEvent) When() time.Time {
	return e.when
}

func (e *banEvent) Why() string {
	return e.why
}

func (e *banEvent) String() string {
	return fmt.Sprintf("%s on %v",e.why,e.when)
}

type byTime []BanEvent
func (a byTime) Len() int           { return len(a) }
func (a byTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byTime) Less(i, j int) bool { return a[i].When().Before(a[j].When()) }

