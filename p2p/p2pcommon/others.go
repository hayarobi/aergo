package p2pcommon

import (
	"io"
	"time"

	"github.com/aergoio/aergo-actor/actor"
	"github.com/aergoio/aergo/message"
	"github.com/aergoio/aergo/types"
	host "github.com/libp2p/go-libp2p-host"
	inet "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	protocol "github.com/libp2p/go-libp2p-protocol"
)

// PeerAccessor is an interface for a another actor module to get info of peers
type PeerAccessor interface {
	GetPeerBlockInfos() []types.PeerBlockInfo
	GetPeer(ID peer.ID) (RemotePeer, bool)
}

// MsgOrder is abstraction of information about the message that will be sent to peer.
// Some type of msgOrder, such as notice mo, should thread-safe and re-entrant
type MsgOrder interface {
	GetMsgID() MsgID
	// Timestamp is unit time value
	Timestamp() int64
	IsRequest() bool
	IsNeedSign() bool
	GetProtocolID() SubProtocol

	// SendTo send message to remote peer. it return err if write fails, or nil if write is successful or ignored.
	SendTo(p RemotePeer) error
}

// ResponseReceiver returns true when receiver handled it, or false if this receiver is not the expected handler.
// NOTE: the return value is temporal works for old implementation and will be remove later.
type MoFactory interface {
	NewMsgRequestOrder(expecteResponse bool, protocolID SubProtocol, message MessageBody) MsgOrder
	NewMsgBlockRequestOrder(respReceiver ResponseReceiver, protocolID SubProtocol, message MessageBody) MsgOrder
	NewMsgResponseOrder(reqID MsgID, protocolID SubProtocol, message MessageBody) MsgOrder
	NewMsgBlkBroadcastOrder(noticeMsg *types.NewBlockNotice) MsgOrder
	NewMsgTxBroadcastOrder(noticeMsg *types.NewTransactionsNotice) MsgOrder
	NewMsgBPBroadcastOrder(noticeMsg *types.BlockProducedNotice) MsgOrder
}

// PeerManager is internal service that provide peer management
type PeerManager interface {
	Start() error
	Stop() error

	//NetworkTransport
	SelfMeta() PeerMeta
	SelfNodeID() peer.ID

	AddNewPeer(peer PeerMeta)
	// Remove peer from peer list. Peer dispose relative resources and stop itself, and then call RemovePeer to peermanager
	RemovePeer(peer RemotePeer)

	NotifyPeerAddressReceived([]PeerMeta)

	// GetPeer return registered(handshaked) remote peer object
	GetPeer(ID peer.ID) (RemotePeer, bool)
	GetPeers() []RemotePeer
	GetPeerAddresses(noHidden bool, showSelf bool) []*message.PeerInfo

	GetPeerBlockInfos() []types.PeerBlockInfo
}
type SyncManager interface {
	// handle notice from bp
	HandleBlockProducedNotice(peer RemotePeer, block *types.Block)
	// handle notice from other node
	HandleNewBlockNotice(peer RemotePeer, data *types.NewBlockNotice)
	HandleGetBlockResponse(peer RemotePeer, msg Message, resp *types.GetBlockResponse)
	HandleNewTxNotice(peer RemotePeer, hashes []types.TxID, data *types.NewTransactionsNotice)
}

// ActorService is collection of helper methods to use actor
// FIXME move to more general package. it used in p2p and rpc
type ActorService interface {
	// TellRequest send actor request, which does not need to get return value, and forget it.
	TellRequest(actor string, msg interface{})
	// SendRequest send actor request, and the response is expected to go back asynchronously.
	SendRequest(actor string, msg interface{})
	// CallRequest send actor request and wait the handling of that message to finished,
	// and get return value.
	CallRequest(actor string, msg interface{}, timeout time.Duration) (interface{}, error)
	// CallRequestDefaultTimeout is CallRequest with default timeout
	CallRequestDefaultTimeout(actor string, msg interface{}) (interface{}, error)

	// FutureRequest send actor reqeust and get the Future object to get the state and return value of message
	FutureRequest(actor string, msg interface{}, timeout time.Duration) *actor.Future
	// FutureRequestDefaultTimeout is FutureRequest with default timeout
	FutureRequestDefaultTimeout(actor string, msg interface{}) *actor.Future

	GetChainAccessor() types.ChainAccessor
}

// NTContainer can provide NetworkTransport interface.
type NTContainer interface {
	GetNetworkTransport() NetworkTransport

	// ChainID return id of current chain.
	ChainID() *types.ChainID
}

// NetworkTransport do manager network connection
// TODO need refactoring. it has other role, pk management of self peer
type NetworkTransport interface {
	host.Host
	Start() error
	Stop() error

	SelfMeta() PeerMeta

	GetAddressesOfPeer(peerID peer.ID) []string

	// AddStreamHandler wrapper function which call host.SetStreamHandler after transport is initialized, this method is for preventing nil error.
	AddStreamHandler(pid protocol.ID, handler inet.StreamHandler)

	GetOrCreateStream(meta PeerMeta, protocolIDs ...protocol.ID) (inet.Stream, error)
	GetOrCreateStreamWithTTL(meta PeerMeta, ttl time.Duration, protocolIDs ...protocol.ID) (inet.Stream, error)

	FindPeer(peerID peer.ID) bool
	ClosePeerConnection(peerID peer.ID) bool
}

// FlushableWriter is writer which have Flush method, such as bufio.Writer
type FlushableWriter interface {
	io.Writer
	// Flush writes any buffered data to the underlying io.Writer.
	Flush() error
}