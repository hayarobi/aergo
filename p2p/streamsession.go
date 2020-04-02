package p2p

import (
	"fmt"
	"github.com/aergoio/aergo-lib/log"
	"github.com/aergoio/aergo/p2p/p2pcommon"
	"github.com/aergoio/aergo/p2p/p2putil"
	"github.com/aergoio/aergo/types"
	"runtime/debug"
	"time"
)

type sessionType int
const (
	baseP2P = iota
	raft
)

type msgHandler interface {
	p2pcommon.RemotePeer

	handleMsg(msg p2pcommon.Message) error
	ReportFail(err error)
}


type StreamSession interface {
	Start()
	Stop(goAwayMsg *types.GoAwayNotice)
	// SendMessage queue msg to send. it return true if queueing is successful, or return false if queue is full
	SendMessage(msg p2pcommon.MsgOrder) bool
	SendAndWaitMessage(msg p2pcommon.MsgOrder, timeout time.Duration) bool

	writeToPeer(m p2pcommon.MsgOrder)
}

type streamSession struct {
	peer msgHandler
	rw   p2pcommon.MsgReadWriter
	st sessionType

	logger *log.Logger

	writeBuf   chan p2pcommon.MsgOrder
	closeWrite chan *types.GoAwayNotice
}

func NewStreamSession(peer p2pcommon.RemotePeer, rw p2pcommon.MsgReadWriter, logger *log.Logger, st sessionType) *streamSession {
	ss := &streamSession{peer:peer.(*remotePeerImpl), rw:rw, st:st, logger:logger,
		writeBuf: make(chan p2pcommon.MsgOrder, writeMsgBufferSize),
		closeWrite: make(chan *types.GoAwayNotice),
	}
	return ss
}

func (ss *streamSession) Start() {
	go ss.runWrite()
	go ss.runRead()
}

func (ss *streamSession) Stop(goAwayMsg *types.GoAwayNotice) {
	ss.closeWrite <- goAwayMsg
}

func (ss *streamSession) SendMessage(msg p2pcommon.MsgOrder) bool {
	select {
	case ss.writeBuf <- msg:
		// it's OK
		return true
	default:
		return false
	}
}

func (ss *streamSession) SendAndWaitMessage(msg p2pcommon.MsgOrder, timeout time.Duration) bool {
	select {
	case ss.writeBuf <- msg:
		// it's OK
		return true
	case <-time.NewTimer(timeout).C:
		return false
	}
}
func (ss *streamSession) runWrite() {
	defer func() {
		if r := recover(); r != nil {
			ss.logger.Panic().Str("callStack", string(debug.Stack())).Str(p2putil.LogPeerName, ss.peer.Name()).Str("recover", fmt.Sprint(r)).Msg("There were panic in runWrite ")
		}
	}()

WRITELOOP:
	for {
		select {
		case m := <-ss.writeBuf:
			ss.writeToPeer(m)
		case lastWord := <-ss.closeWrite:
			l := ss.logger.Debug().Str(p2putil.LogPeerName, ss.peer.Name())
			if lastWord != nil {
				l.Str("goAwayMsg",lastWord.Message)
			}
			l.Msg("Quitting runWrite")
			close(ss.closeWrite)
			if lastWord != nil {
				msgVal := p2pcommon.NewSimpleMsgVal(p2pcommon.GoAway, p2pcommon.EmptyID)
				bytes, _ := p2putil.MarshalMessageBody(lastWord)
				msgVal.SetPayload(bytes)
				ss.rw.WriteMsg(msgVal)
			}
			break WRITELOOP
		}
	}
	ss.cleanupWrite()
	// how to close more gracefully?
	ss.rw.Close()
}

func (ss *streamSession) cleanupWrite() {
	// 1. cleaning receive handlers. TODO add code

	// 2. canceling not sent orders
	for {
		select {
		case m := <-ss.writeBuf:
			m.CancelSend(nil)
		default:
			return
		}
	}
}

func (ss *streamSession) runRead() {
	for {
		msg, err := ss.rw.ReadMsg()
		if err != nil {
			select {
			case <- ss.closeWrite:
				// maybe local peer closed this session
				ss.logger.Trace().Str(p2putil.LogPeerName, ss.peer.Name()).Err(err).Msg("closing read from finished session")
			default:
					// set different log level by case (i.e. it can be expected if peer is disconnecting )
					ss.logger.Warn().Str(p2putil.LogPeerName, ss.peer.Name()).Err(err).Msg("Failed to read message")
					ss.peer.ReportFail(err)
			}
			return
		} else {
			ss.peer.handleMsg(msg)
		}
	}
}

func (ss *streamSession) writeToPeer(m p2pcommon.MsgOrder) {
	if err := m.SendTo(ss.peer, ss.rw); err != nil {
		// write fail
		ss.peer.ReportFail(err)
	}
}
