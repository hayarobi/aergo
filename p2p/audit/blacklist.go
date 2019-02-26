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

	//AddBanScore(addr string, pid peer.ID, why string)

	IsBanned(addr string, pid peer.ID) (bool, time.Time)
	IsBannedPeerID(peerID peer.ID) (bool, time.Time)
	IsBannedAddr(addr string) (bool, time.Time)

	// GetStatusByID returns ban status of peer. it returns error when the peer is not registered to ban
	GetStatusByID(peerID peer.ID) (BanStatus, error)
	// GetStatusByID returns ban status of ip address. it returns error when the peer is not registered to ban
	GetStatusByAddr(addr string) (BanStatus, error)
}

// BanStatus keep kickout logs and decide how long the ban duration is
type BanStatus interface {
	// ID is ip address or peer id
	ID() string

	// ValidUntil show when this ban items is expired.
	ValidUntil() time.Time
	Events() []BanEvent
}

type BanEvent interface {
	When() time.Time
	Why() string
}

// BanThreshold is number of events to ban address or peerid
const BanThreshold = 5

var BanDurations = []time.Duration{
	0,
	0,
	time.Minute,
	time.Minute*3,
	time.Minute*10,
	time.Hour,
	time.Hour*24,
	time.Hour*24*30,
	time.Hour*24*3650,
}
const BanValidDuration = time.Minute * 30
const BanReleaseDuration = time.Hour * 24 * 730