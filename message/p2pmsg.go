/**
 *  @file
 *  @copyright defined in aergo/LICENSE.txt
 */

package message

import (
	"fmt"
	"github.com/aergoio/aergo/types"
	"github.com/libp2p/go-libp2p-peer"
	"time"
)

const P2PSvc = "p2pSvc"

// errors which async responses of p2p actor, such as GetBlockChunksRsp, can contains,
var (
	RemotePeerFailError  = fmt.Errorf("remote peer return error")
	PeerNotFoundError    = fmt.Errorf("remote peer was not found")
	MissingHashError     = fmt.Errorf("some block hash not found")
	UnexpectedBlockError = fmt.Errorf("unexpected blocks response")
	TooFewBlocksError    = fmt.Errorf("too few blocks received that expected")
	TooManyBlocksError   = fmt.Errorf("too many blocks received that expected")
	TooBigBlockError     = fmt.Errorf("block size limit exceeded")
)

// PingMsg send types.Ping to each peer.
// The actor returns true if sending is successful.
type PingMsg struct {
	ToWhom peer.ID
}

// GetAddressesMsg send types.AddressesRequest to dest peer. the dest peer will send types.AddressesResponse.
// The actor returns true if sending is successful.
type GetAddressesMsg struct {
	ToWhom peer.ID
	Size   uint32
	Offset uint32
}

// NotifyNewBlock send types.NewBlockNotice to other peers. The receiving peer will send GetBlockHeadersRequest or GetBlockRequest if needed.
// The actor returns true if sending is successful.
type NotifyNewBlock struct {
	Produced bool
	BlockNo  uint64
	Block    *types.Block
}

type BlockHash []byte
type TXHash []byte

// NotifyNewTransactions send types.NewTransactionsNotice to other peers.
// The actor returns true if sending is successful.
type NotifyNewTransactions struct {
	Txs []*types.Tx
}

// GetTransactions send types.GetTransactionsRequest to dest peer. The receiving peer will send types.GetTransactionsResponse
// The actor returns true if sending is successful.
type GetTransactions struct {
	ToWhom peer.ID
	Hashes []TXHash
}

// TransactionsResponse is data from other peer, as a response of types.GetTransactionsRequest
// p2p module will send this to mempool actor.
type TransactionsResponse struct {
	txs []*types.Tx
}

// GetBlockHeaders send type.GetBlockRequest to dest peer
// The actor returns true if sending is successful.
type GetBlockHeaders struct {
	ToWhom peer.ID
	// Hash is the first block to get. Height will be used when Hash mi empty
	Hash    BlockHash
	Height  uint64
	Asc     bool
	Offset  uint64
	MaxSize uint32
}

// BlockHeadersResponse is data from other peer, as a response of types.GetBlockRequest
// p2p module will send this to chainservice actor.
type BlockHeadersResponse struct {
	Hashes  []BlockHash
	Headers []*types.BlockHeader
}

// GetBlockInfos send types.GetBlockRequest to dest peer.
// The actor returns true if sending is successful.
type GetBlockInfos struct {
	ToWhom peer.ID
	Hashes []BlockHash
}

type GetBlockChunks struct {
	Seq uint64
	GetBlockInfos
	TTL time.Duration
}

// BlockInfosResponse is data from other peer, as a response of types.GetBlockRequest
// p2p module will send this to chainservice actor.
type BlockInfosResponse struct {
	FromWhom peer.ID
	Blocks   []*types.Block
}

type GetBlockChunksRsp struct {
	Seq    uint64
	ToWhom peer.ID
	Blocks []*types.Block
	Err    error
}

// GetPeers requests p2p actor to get remote peers that is connected.
// The actor returns *GetPeersRsp
type GetPeers struct {
	NoHidden bool
	ShowSelf bool
}

type PeerInfo struct {
	Addr            *types.PeerAddress
	Version         string
	Hidden          bool
	CheckTime       time.Time
	LastBlockHash   []byte
	LastBlockNumber uint64
	State           types.PeerState
	Self            bool
}

// GetPeersRsp contains peer meta information and current states.
type GetPeersRsp struct {
	Peers []*PeerInfo
}

type GetMetrics struct {
}

// GetSyncAncestor is sent from Syncer, send types.GetAncestorRequest to dest peer.
type GetSyncAncestor struct {
	Seq    uint64
	ToWhom peer.ID
	Hashes [][]byte
}

// GetSyncAncestorRsp is data from other peer, as a response of types.GetAncestorRequest
type GetSyncAncestorRsp struct {
	Seq      uint64
	Ancestor *types.BlockInfo
}

type GetHashes struct {
	Seq      uint64
	ToWhom   peer.ID
	PrevInfo *types.BlockInfo
	Count    uint64
}

type GetHashesRsp struct {
	Seq      uint64
	PrevInfo *types.BlockInfo
	Hashes   []BlockHash
	Count    uint64
	Err      error
}

type GetHashByNo struct {
	Seq     uint64
	ToWhom  peer.ID
	BlockNo types.BlockNo
}

type GetHashByNoRsp struct {
	Seq       uint64
	BlockHash BlockHash
	Err       error
}

type GetSelf struct {
}

type GetCluster struct {
	ReplyC chan *GetClusterRsp
}

type GetClusterRsp struct {
	ChainID BlockHash
	Members []*types.MemberAttr
	Err     error
}
// SenderContext is additional information about a remote peer causing inter-actor message
type SenderContext struct {
	PeerID       peer.ID
	ManageNumber uint32
}

// BlamableError is error which blames message sender for bad request, and contains how big the sender's fault is
type BlamableError interface {
	error
	Size() FaultSize
}


// FaultSize is how big the sender's fault
type FaultSize int32

const (
	Tiny FaultSize = iota
	Small
	Normal
	Big
	// Severe will disconnect peer immediately and add to blacklist to prevent further connection attempt
	Severe
)

type StringBlameError struct {
	size  FaultSize
	errmsg string
}

func NewBlamableError(faultSize FaultSize, errmsg string) BlamableError {
	return &StringBlameError{size: faultSize, errmsg:errmsg}
}

func (e *StringBlameError) Error() string {
	return e.errmsg
}

func (e *StringBlameError) Size() FaultSize {
	return e.size
}

type BlamableWrapper struct {
	size  FaultSize
	inner error
}

func NewBlamableWrapper(faultSize FaultSize, originalError error) BlamableError {
	return &BlamableWrapper{size: faultSize, inner:originalError}
}

func (bw *BlamableWrapper) Size() FaultSize {
	return bw.size
}

func (bw *BlamableWrapper) Error() string {
	return bw.inner.Error()
}


