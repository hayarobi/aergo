/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package subproto

import (
	"github.com/aergoio/aergo-lib/log"
	"github.com/aergoio/aergo/internal/enc"
	"github.com/aergoio/aergo/message"
	"github.com/aergoio/aergo/p2p/p2pcommon"
	"github.com/aergoio/aergo/p2p/p2putil"
	"github.com/aergoio/aergo/types"
	"github.com/golang/protobuf/proto"
)

type txRequestHandler struct {
	BaseMsgHandler
	msgHelper message.Helper
}

var _ p2pcommon.MessageHandler = (*txRequestHandler)(nil)

type txResponseHandler struct {
	BaseMsgHandler
}

var _ p2pcommon.MessageHandler = (*txResponseHandler)(nil)

type newTxNoticeHandler struct {
	BaseMsgHandler
}

var _ p2pcommon.MessageHandler = (*newTxNoticeHandler)(nil)

// newTxReqHandler creates handler for GetTransactionsRequest
func NewTxReqHandler(pm p2pcommon.PeerManager, peer p2pcommon.RemotePeer, logger *log.Logger, actor p2pcommon.ActorService) *txRequestHandler {
	th := &txRequestHandler{
		BaseMsgHandler{protocol: GetTXsRequest, pm: pm, peer: peer, actor: actor, logger: logger},
		message.GetHelper()}
	return th
}

func (th *txRequestHandler) ParsePayload(rawbytes []byte) (p2pcommon.MessageBody, error) {
	return p2putil.UnmarshalAndReturn(rawbytes, &types.GetTransactionsRequest{})
}

func (th *txRequestHandler) Handle(msg p2pcommon.Message, msgBody p2pcommon.MessageBody) {

	remotePeer := th.peer
	reqHashes := msgBody.(*types.GetTransactionsRequest).Hashes
	p2putil.DebugLogReceiveMsg(th.logger, th.protocol, msg.ID().String(), remotePeer, p2putil.BytesArrToString(reqHashes))

	// TODO consider to make async if deadlock with remote peer can occurs
	// NOTE size estimation is tied to protobuf3 it should be changed when protobuf is changed.
	// find transactions from chainservice
	idx := 0
	status := types.ResultStatus_OK
	var hashes []types.TxHash
	var txInfos, txs []*types.Tx
	payloadSize := EmptyGetBlockResponseSize
	var txSize, fieldSize int

	bucket := message.MaxReqestHashes
	var futures []interface{}

	for _, h := range reqHashes {
		hashes = append(hashes, h)
		if len(hashes) == bucket {
			if f, err := th.actor.CallRequestDefaultTimeout(message.MemPoolSvc,
				&message.MemPoolExistEx{Hashes: hashes}); err == nil {
				futures = append(futures, f)
			}
			hashes = nil
		}
	}
	if hashes != nil {
		if f, err := th.actor.CallRequestDefaultTimeout(message.MemPoolSvc,
			&message.MemPoolExistEx{Hashes: hashes}); err == nil {
			futures = append(futures, f)
		}
	}
	hashes = nil
	idx = 0
	for _, f := range futures {
		if tmp, err := th.msgHelper.ExtractTxsFromResponseAndError(f, nil); err == nil {
			txs = append(txs, tmp...)
		} else {
			th.logger.Debug().Err(err).Msg("ErrExtract tx in future")
		}
	}
	for _, tx := range txs {
		if tx == nil {
			continue
		}
		hash := tx.GetHash()
		txSize = proto.Size(tx)

		fieldSize = txSize + p2putil.CalculateFieldDescSize(txSize)
		fieldSize += len(hash) + p2putil.CalculateFieldDescSize(len(hash))

		if (payloadSize + fieldSize) > p2pcommon.MaxPayloadLength {
			// send partial list
			resp := &types.GetTransactionsResponse{
				Status: status,
				Hashes: hashes,
				Txs:    txInfos, HasNext: true}
			th.logger.Debug().Int(p2putil.LogTxCount, len(hashes)).
				Str(p2putil.LogOrgReqID, msg.ID().String()).Msg("Sending partial response")

			remotePeer.SendMessage(remotePeer.MF().
				NewMsgResponseOrder(msg.ID(), GetTXsResponse, resp))
			hashes, txInfos, payloadSize = nil, nil, EmptyGetBlockResponseSize
		}

		hashes = append(hashes, hash)
		txInfos = append(txInfos, tx)
		payloadSize += fieldSize
		idx++
	}
	if 0 == idx {
		status = types.ResultStatus_NOT_FOUND
	}
	th.logger.Debug().Int(p2putil.LogTxCount, len(hashes)).
		Str(p2putil.LogOrgReqID, msg.ID().String()).Str(p2putil.LogRespStatus, status.String()).Msg("Sending last part response")
	// generate response message

	resp := &types.GetTransactionsResponse{
		Status: status,
		Hashes: hashes,
		Txs:    txInfos, HasNext: false}
	remotePeer.SendMessage(remotePeer.MF().NewMsgResponseOrder(msg.ID(), GetTXsResponse, resp))
}

// newTxRespHandler creates handler for GetTransactionsResponse
func NewTxRespHandler(pm p2pcommon.PeerManager, peer p2pcommon.RemotePeer, logger *log.Logger, actor p2pcommon.ActorService) *txResponseHandler {
	th := &txResponseHandler{BaseMsgHandler: BaseMsgHandler{protocol: GetTXsResponse, pm: pm, peer: peer, actor: actor, logger: logger}}
	return th
}

func (th *txResponseHandler) ParsePayload(rawbytes []byte) (p2pcommon.MessageBody, error) {
	return p2putil.UnmarshalAndReturn(rawbytes, &types.GetTransactionsResponse{})
}

func (th *txResponseHandler) Handle(msg p2pcommon.Message, msgBody p2pcommon.MessageBody) {
	peer := th.peer
	data := msgBody.(*types.GetTransactionsResponse)
	p2putil.DebugLogReceiveResponseMsg(th.logger, th.protocol, msg.ID().String(), msg.OriginalID().String(), th.peer, len(data.Txs))

	th.peer.ConsumeRequest(msg.OriginalID())
	// TODO: Is there any better solution than passing everything to mempool service?
	if len(data.Txs) > 0 {
		th.logger.Debug().Int(p2putil.LogTxCount, len(data.Txs)).Msg("Request mempool to add txs")
		//th.actor.SendRequest(message.MemPoolSvc, &message.MemPoolPut{Txs: data.Txs})
		for _, tx := range data.Txs {
			th.actor.SendRequest(message.MemPoolSvc, &message.MemPoolPut{Tx: tx, Sender:&message.SenderContext{peer.ID(), peer.ManageNumber()}})
		}
	}
}

// newNewTxNoticeHandler creates handler for GetTransactionsResponse
func NewNewTxNoticeHandler(pm p2pcommon.PeerManager, peer p2pcommon.RemotePeer, logger *log.Logger, actor p2pcommon.ActorService, sm p2pcommon.SyncManager) *newTxNoticeHandler {
	th := &newTxNoticeHandler{BaseMsgHandler: BaseMsgHandler{protocol: NewTxNotice, pm: pm, sm: sm, peer: peer, actor: actor, logger: logger}}
	return th
}

func (th *newTxNoticeHandler) ParsePayload(rawbytes []byte) (p2pcommon.MessageBody, error) {
	return p2putil.UnmarshalAndReturn(rawbytes, &types.NewTransactionsNotice{})
}

func (th *newTxNoticeHandler) Handle(msg p2pcommon.Message, msgBody p2pcommon.MessageBody) {
	remotePeer := th.peer
	data := msgBody.(*types.NewTransactionsNotice)
	// remove to verbose log
	if th.logger.IsDebugEnabled() {
		p2putil.DebugLogReceiveMsg(th.logger, th.protocol, msg.ID().String(), remotePeer, p2putil.BytesArrToString(data.TxHashes))
	}

	if len(data.TxHashes) == 0 {
		return
	}
	// lru cache can accept hashable key
	hashes := make([]types.TxID, len(data.TxHashes))
	for i, hash := range data.TxHashes {
		if tid, err := types.ParseToTxID(hash); err != nil {
			th.logger.Info().Str(p2putil.LogPeerName, remotePeer.Name()).Str("hash", enc.ToString(hash)).Msg("malformed txhash found")
			// TODO Add penelty score and break
			break
		} else {
			hashes[i] = tid
		}
	}
	added := th.peer.UpdateTxCache(hashes)
	if len(added) > 0 {
		th.sm.HandleNewTxNotice(th.peer, added, data)
	}
}
