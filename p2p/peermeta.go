/**
 *  @file
 *  @copyright defined in aergo/LICENSE.txt
 */

package p2p

import (
	"fmt"
	"time"

	"strconv"

	"github.com/aergoio/aergo/types"
	"github.com/libp2p/go-libp2p-peer"
)

const (
	UnknownName = "unknown"
	NotManaged = 0
)
// PeerMeta contains non changeable information of peer node during connected state
// TODO: PeerMeta is almost same as PeerAddress, so TODO to unify them.
type PeerMeta struct {
	ID peer.ID
	Nick string

	// IPAddress is human readable form of ip address such as "192.168.0.1" or "2001:0db8:0a0b:12f0:33:1"
	IPAddress  string
	Port       uint32
	Designated bool // Designated means this peer is designated in config file and connect to in startup phase

	Hidden    bool // Hidden means that meta info of this peer will not be sent to other peers when getting peer list
	Outbound   bool
}

func (m *PeerMeta) Pretty() string {
	pid := m.ID.Pretty()
	if idlen := len(pid); idlen > 10 {
	 	return fmt.Sprintf("%s:%s..%s",m.Nick,pid[:2],pid[len(pid)-6:])
	} else {
		return fmt.Sprintf("%s:%s",m.Nick,pid)
	}
}

func (m PeerMeta) String() string {
	return m.ID.Pretty() + "/" + m.IPAddress + ":" + strconv.Itoa(int(m.Port))
}

// FromPeerAddress convert PeerAddress to PeerMeta
func FromPeerAddress(addr *types.PeerAddress) PeerMeta {
	meta := PeerMeta{IPAddress: addr.Address, Nick:UnknownName,
		Port: addr.Port, ID: peer.ID(addr.PeerID)}
	return meta
}

// ToPeerAddress convert PeerMeta to PeerAddress
func (m PeerMeta) ToPeerAddress() types.PeerAddress {
	addr := types.PeerAddress{Address: m.IPAddress, Port: m.Port,
		PeerID: []byte(m.ID)}
	return addr
}

// TTL return node's ttl
func (m PeerMeta) TTL() time.Duration {
	if m.Designated {
		return DesignatedNodeTTL
	}
	return DefaultNodeTTL
}
