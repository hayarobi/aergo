/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package audit

import (
	"github.com/libp2p/go-libp2p-peer"
	"time"
)

type banStatusImpl struct {
	events []BanEvent
	banScore int

	banUntil time.Time
}

func newStatus(initialScore int, banTime time.Time, prevEvents []BanEvent) *banStatusImpl {
	tmp := make([]BanEvent, len(prevEvents))
	copy(tmp, prevEvents)

	return &banStatusImpl{banScore:initialScore, events:tmp, banUntil:banTime}
}

func (bs *banStatusImpl) ValidUntil() time.Time {
	return bs.banUntil
}

func (bs *banStatusImpl) Events() []BanEvent {
	return bs.events
}

func (bs *banStatusImpl)addEvent(ev BanEvent) {
	// this method should be inside lock
	// recalculate total score
	// add last event
	bs.events = append(bs.events, ev)
	bs.banScore++

	bs.updateStats(time.Now())
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

func newIDBanStatusImpl() *idBanStatusImpl {
	return &idBanStatusImpl{banStatusImpl:banStatusImpl{banUntil: UndefinedTime}}
}


func (bs *idBanStatusImpl) ID() string {
	return bs.id.Pretty()
}


type addrBanStatusImpl struct {
	banStatusImpl
	addr string
}

func newAddrBanStatusImpl() *addrBanStatusImpl {
	return &addrBanStatusImpl{banStatusImpl:banStatusImpl{banUntil: UndefinedTime}}
}

func (bs *addrBanStatusImpl) ID() string {
	return bs.addr
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

type byTime []BanEvent
func (a byTime) Len() int           { return len(a) }
func (a byTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byTime) Less(i, j int) bool { return a[i].When().Before(a[j].When()) }

