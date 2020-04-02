package p2p

import (
	"github.com/aergoio/aergo/p2p/p2pcommon"
	"github.com/aergoio/aergo/p2p/p2pmock"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)


func TestRemotePeer_writeToPeer(t *testing.T) {
	rand := uuid.Must(uuid.NewV4())
	var sampleMsgID p2pcommon.MsgID
	copy(sampleMsgID[:], rand[:])
	type args struct {
		StreamResult error
		signErr      error
		needResponse bool
		sendErr      error
	}
	type wants struct {
		sendCnt   int
		expReqCnt int
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{"TNReq1", args{}, wants{1, 0}},
		// {"TNReqWResp1", args{needResponse: true}, wants{1, 1}},

		// error while signing
		// error while get stream
		// TODO this case is evaluated in pbMsgOrder tests
		// {"TFSend1", args{needResponse: true, sendErr: sampleErr}, wants{1, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			dummyPeer := &remotePeerImpl{}
			mockMO := p2pmock.NewMockMsgOrder(ctrl)
			mockStream := p2pmock.NewMockStream(ctrl)
			dummyRW := p2pmock.NewMockMsgReadWriter(ctrl)
			dummyRW.EXPECT().Close().AnyTimes()

			mockStream.EXPECT().Close().Return(nil).AnyTimes()
			mockMO.EXPECT().IsNeedSign().Return(true).AnyTimes()
			mockMO.EXPECT().SendTo(gomock.Any()).Return(tt.args.sendErr)
			mockMO.EXPECT().GetProtocolID().Return(p2pcommon.PingRequest).AnyTimes()
			mockMO.EXPECT().GetMsgID().Return(sampleMsgID).AnyTimes()

			//sampleConn := p2pcommon.RemoteConn{IP:net.ParseIP(sampleMeta.PrimaryAddress()),Port:sampleMeta.PrimaryPort()}
			//sampleRemote := p2pcommon.RemoteInfo{Meta:sampleMeta, Connection:sampleConn}

			ss := NewStreamSession(dummyPeer, dummyRW, logger, baseP2P )
			go ss.runWrite()
			ss.writeToPeer(mockMO)

			// FIXME wait in more reliable way
			time.Sleep(50 * time.Millisecond)
			ss.closeWrite <- nil
			//mockOrder.AssertNumberOfCalls(t, "SendTo", tt.wants.sendCnt)
			//assert.Equal(t, tt.wants.expReqCnt, len(p.requests))
		})
	}
}
