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
}

func (*banStatusImpl) ExpireAt() time.Time {
	panic("implement me")
}

func (bs *banStatusImpl) Events() []BanEvent {
	return bs.events
}

func (bs *banStatusImpl)addEvent(event BanEvent) {
	bs.events = append(bs.events, event)
}

type idBanStatusImpl struct {
	banStatusImpl
	id peer.ID
}

func newIDBanStatusImpl() *idBanStatusImpl {
	return &idBanStatusImpl{}
}


func (bs *idBanStatusImpl) ID() string {
	return bs.id.Pretty()
}


type addrBanStatusImpl struct {
	banStatusImpl
	addr string
}

func newAddrBanStatusImpl() *addrBanStatusImpl {
	return &addrBanStatusImpl{}
}

func (bs *addrBanStatusImpl) ID() string {
	return bs.addr
}



type banEvent struct {
	when time.Time
	why string
}

func (e banEvent) When() time.Time {
	return e.when
}

func (e banEvent) Why() string {
	return e.why
}
