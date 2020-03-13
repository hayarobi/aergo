/**
 *  @file
 *  @copyright defined in aergo/LICENSE.txt
 */

package p2p

import (
	"github.com/aergoio/aergo-lib/log"
	"github.com/aergoio/aergo/consensus"
	"github.com/aergoio/aergo/p2p/raftsupport"
	"github.com/aergoio/etcd/raft/raftpb"
	"time"

	"github.com/aergoio/aergo/internal/enc"
	"github.com/aergoio/aergo/p2p/p2pcommon"
	"github.com/aergoio/aergo/p2p/p2putil"

	"github.com/aergoio/aergo/types"
)

// ClientVersion is the version of p2p protocol to which this codes are built

type pbMessageOrder struct {
	logger *log.Logger
	request    bool
	needSign   bool
	trace      bool
	protocolID p2pcommon.SubProtocol // protocolName and msg struct type MUST be matched.

	message p2pcommon.Message
}

var _ p2pcommon.MsgOrder = (*pbRequestOrder)(nil)
var _ p2pcommon.MsgOrder = (*pbResponseOrder)(nil)
var _ p2pcommon.MsgOrder = (*pbBlkNoticeOrder)(nil)
var _ p2pcommon.MsgOrder = (*pbTxNoticeOrder)(nil)
var _ p2pcommon.MsgOrder = (*pbRaftMsgOrder)(nil)

func (mo *pbMessageOrder) GetMsgID() p2pcommon.MsgID {
	return mo.message.ID()
}

func (mo *pbMessageOrder) Timestamp() int64 {
	return mo.message.Timestamp()
}

func (mo *pbMessageOrder) IsRequest() bool {
	return mo.request
}

func (mo *pbMessageOrder) IsNeedSign() bool {
	return mo.needSign
}

func (mo *pbMessageOrder) GetProtocolID() p2pcommon.SubProtocol {
	return mo.protocolID
}

func (mo *pbMessageOrder) CancelSend(pi p2pcommon.RemotePeer) {
}

type pbRequestOrder struct {
	pbMessageOrder
	respReceiver p2pcommon.ResponseReceiver
}

func (rmo *pbRequestOrder) SendTo(p p2pcommon.RemotePeer, rw p2pcommon.MsgReadWriter) error {
	p.RegisterRequest(&p2pcommon.RequestInfo{CTime: time.Now(), ReqMO: rmo, Receiver: rmo.respReceiver})

	err := rw.WriteMsg(rmo.message)
	if err != nil {
		rmo.logger.Warn().Str(p2putil.LogPeerName, p.Name()).Str(p2putil.LogProtoID, rmo.GetProtocolID().String()).Str(p2putil.LogMsgID, rmo.GetMsgID().String()).Err(err).Msg("fail to SendTo")
		p.ConsumeRequest(rmo.GetMsgID())
		return err
	}

	if rmo.trace {
		rmo.logger.Debug().Str(p2putil.LogPeerName, p.Name()).Str(p2putil.LogProtoID, rmo.GetProtocolID().String()).
			Str(p2putil.LogMsgID, rmo.GetMsgID().String()).Msg("Send request message")
	}
	return nil
}

type pbResponseOrder struct {
	pbMessageOrder
}

func (rmo *pbResponseOrder) SendTo(p p2pcommon.RemotePeer, rw p2pcommon.MsgReadWriter) error {
	err := rw.WriteMsg(rmo.message)
	if err != nil {
		rmo.logger.Warn().Str(p2putil.LogPeerName, p.Name()).Str(p2putil.LogProtoID, rmo.GetProtocolID().String()).Str(p2putil.LogMsgID, rmo.GetMsgID().String()).Err(err).Msg("fail to SendTo")
		return err
	}
	if rmo.trace {
		rmo.logger.Debug().Str(p2putil.LogPeerName, p.Name()).Str(p2putil.LogProtoID, rmo.GetProtocolID().String()).
			Str(p2putil.LogMsgID, rmo.GetMsgID().String()).Str(p2putil.LogOrgReqID, rmo.message.OriginalID().String()).Msg("Send response message")
	}

	return nil
}

type pbBlkNoticeOrder struct {
	pbMessageOrder
	blkHash []byte
	blkNo   uint64
}

func (rmo *pbBlkNoticeOrder) SendTo(pi p2pcommon.RemotePeer, rw p2pcommon.MsgReadWriter) error {
	p := pi.(*remotePeerImpl)
	var blkHash = types.ToBlockID(rmo.blkHash)
	passedTime := time.Now().Sub(p.lastBlkNoticeTime)
	skipNotice := false
	if p.LastStatus().BlockNumber >= rmo.blkNo {
		heightDiff := p.LastStatus().BlockNumber - rmo.blkNo
		switch {
		case heightDiff >= GapToSkipAll:
			skipNotice = true
		case heightDiff >= GapToSkipHourly:
			skipNotice = p.skipCnt < GapToSkipHourly
		default:
			skipNotice = p.skipCnt < GapToSkip5Min
		}
	}
	if skipNotice || passedTime < MinNewBlkNoticeInterval {
		p.skipCnt++
		if p.skipCnt&0x03ff == 0 && rmo.trace {
			rmo.logger.Debug().Str(p2putil.LogPeerName, p.Name()).Str(p2putil.LogProtoID, rmo.GetProtocolID().String()).Int32("skip_cnt", p.skipCnt).Msg("Skipped NewBlockNotice ")
		}
		return nil
	}

	if ok, _ := p.blkHashCache.ContainsOrAdd(blkHash, cachePlaceHolder); ok {
		// the remote peer already know this block hash. skip too many not-interesting log,
		// rmo.logger.Debug().Str(LogPeerName,p.Name()).Str(LogProtoID, rmo.GetProtocolID().String()).
		// 	Str(LogMsgID, rmo.GetMsgID()).Msg("Cancel sending blk notice. peer knows this block")
		return nil
	}
	err := rw.WriteMsg(rmo.message)
	if err != nil {
		rmo.logger.Warn().Str(p2putil.LogPeerName, p.Name()).Str(p2putil.LogProtoID, rmo.GetProtocolID().String()).Str(p2putil.LogMsgID, rmo.GetMsgID().String()).Err(err).Msg("fail to SendTo")
		return err
	}
	p.lastBlkNoticeTime = time.Now()
	if p.skipCnt > 100 && rmo.trace {
		rmo.logger.Debug().Str(p2putil.LogPeerName, p.Name()).Str(p2putil.LogProtoID, rmo.GetProtocolID().String()).Int32("skip_cnt", p.skipCnt).Msg("Send NewBlockNotice after long skip")
	}
	p.skipCnt = 0
	return nil
}

type pbBpNoticeOrder struct {
	pbMessageOrder
	block *types.Block
}

func (rmo *pbBpNoticeOrder) SendTo(pi p2pcommon.RemotePeer, rw p2pcommon.MsgReadWriter) error {
	p := pi.(*remotePeerImpl)
	var blkHash = types.ToBlockID(rmo.block.Hash)
	p.blkHashCache.ContainsOrAdd(blkHash, cachePlaceHolder)
	err := p.rw.WriteMsg(rmo.message)
	if err != nil {
		rmo.logger.Warn().Str(p2putil.LogPeerName, p.Name()).Str(p2putil.LogProtoID, rmo.GetProtocolID().String()).Str(p2putil.LogMsgID, rmo.GetMsgID().String()).Err(err).Msg("fail to SendTo")
		return err
	}
	if rmo.trace {
		rmo.logger.Debug().Str(p2putil.LogPeerName, p.Name()).Str(p2putil.LogProtoID, rmo.GetProtocolID().String()).
			Str(p2putil.LogMsgID, rmo.GetMsgID().String()).Str(p2putil.LogBlkHash, enc.ToString(rmo.block.Hash)).Msg("Notify block produced")
	}
	return nil
}

type pbTxNoticeOrder struct {
	pbMessageOrder
	txHashes []types.TxID
}

func (rmo *pbTxNoticeOrder) SendTo(pi p2pcommon.RemotePeer, rw p2pcommon.MsgReadWriter) error {
	p := pi.(*remotePeerImpl)

	err := rw.WriteMsg(rmo.message)
	if err != nil {
		rmo.logger.Warn().Str(p2putil.LogPeerName, p.Name()).Str(p2putil.LogProtoID, rmo.GetProtocolID().String()).Str(p2putil.LogMsgID, rmo.GetMsgID().String()).Err(err).Msg("fail to SendTo")
		return err
	}
	if rmo.logger.IsDebugEnabled() && rmo.trace {
		rmo.logger.Debug().Str(p2putil.LogPeerName, p.Name()).Str(p2putil.LogProtoID, rmo.GetProtocolID().String()).
			Str(p2putil.LogMsgID, rmo.GetMsgID().String()).Int("hash_cnt", len(rmo.txHashes)).Array("hashes", types.NewLogTxIDsMarshaller(rmo.txHashes, 10)).Msg("Sent tx notice")
	}
	return nil
}

func (rmo *pbTxNoticeOrder) CancelSend(pi p2pcommon.RemotePeer) {
}

type pbRaftMsgOrder struct {
	pbMessageOrder
	raftAcc consensus.AergoRaftAccessor
	msg     *raftpb.Message
}

func (rmo *pbRaftMsgOrder) SendTo(pi p2pcommon.RemotePeer, rw p2pcommon.MsgReadWriter) error {
	p := pi.(*remotePeerImpl)

	err := rw.WriteMsg(rmo.message)
	if err != nil {
		rmo.logger.Warn().Str(p2putil.LogPeerName, p.Name()).Str(p2putil.LogMsgID, rmo.GetMsgID().String()).Err(err).Object("raftMsg", raftsupport.RaftMsgMarshaller{rmo.msg}).Msg("fail to Send raft message")
		rmo.raftAcc.ReportUnreachable(pi.ID())
		return err
	}
	if rmo.trace && rmo.logger.IsDebugEnabled() {
		rmo.logger.Debug().Str(p2putil.LogPeerName, p.Name()).Str(p2putil.LogMsgID, rmo.GetMsgID().String()).Object("raftMsg", raftsupport.RaftMsgMarshaller{rmo.msg}).Msg("Sent raft message")
	}
	return nil
}

func (rmo *pbRaftMsgOrder) CancelSend(pi p2pcommon.RemotePeer) {
	// TODO test more whether to uncomment or to delete code below
	//rmo.raftAcc.ReportUnreachable(pi.ID())
}


type pbTossOrder struct {
	pbMessageOrder
}

func (rmo *pbTossOrder) SendTo(pi p2pcommon.RemotePeer, rw p2pcommon.MsgReadWriter) error {
	p := pi.(*remotePeerImpl)
	err := rw.WriteMsg(rmo.message)
	if err != nil {
		rmo.logger.Warn().Str(p2putil.LogPeerName, p.Name()).Str(p2putil.LogProtoID, rmo.GetProtocolID().String()).Str(p2putil.LogMsgID, rmo.GetMsgID().String()).Err(err).Msg("fail rmo toss")
		return err
	}

	if rmo.trace {
		rmo.logger.Debug().Str(p2putil.LogPeerName, p.Name()).Str(p2putil.LogProtoID, rmo.GetProtocolID().String()).
			Str(p2putil.LogMsgID, rmo.GetMsgID().String()).Msg("toss message")
	}
	return nil
}
