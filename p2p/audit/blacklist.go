/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package audit

import (
	"github.com/libp2p/go-libp2p-peer"
	"time"
)

type BlacklistManager interface {
	NewPeerAuditor(address string, peerID peer.ID, exceedlistener ExceedListener) PeerAuditor

	AddBanScore(addr string, pid peer.ID, why string)

	IsBanned(addr string, pid peer.ID) bool
	IsBannedIP(addr string) bool
	IsBannedAddr(pid string) bool

	// GetStatusByID returns ban status of peer. it returns error when the peer is not registered to ban
	GetStatusByID(peerID peer.ID) (BanStatus, error)
	// GetStatusByID returns ban status of ip address. it returns error when the peer is not registered to ban
	GetStatusByAddr(addr string) (BanStatus, error)
}

type BanStatus interface {
	// ID is ip address or peer id
	ID() string

	// ExpireAt show when this ban items is expired.
	ExpireAt() time.Time
	Events() []BanEvent
}

type BanEvent interface {
	When() time.Time
	Why() string
}
