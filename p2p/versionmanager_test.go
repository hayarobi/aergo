/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package p2p

import (
	"github.com/aergoio/aergo/chain"
	"github.com/aergoio/aergo/p2p/p2pmock"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"

	"github.com/aergoio/aergo/p2p/p2pcommon"
	"github.com/aergoio/aergo/types"
)

func Test_defaultVersionManager_FindBestP2PVersion(t *testing.T) {

	dummyChainID := &types.ChainID{}


	type args struct {
		versions []p2pcommon.P2PVersion
	}
	tests := []struct {
		name   string
		args   args

		want   p2pcommon.P2PVersion
	}{
		{"TCurrent", args{[]p2pcommon.P2PVersion{p2pcommon.P2PVersion230}}, p2pcommon.P2PVersion230},
		{"TLegacy", args{[]p2pcommon.P2PVersion{p2pcommon.P2PVersion200,p2pcommon.P2PVersion033,p2pcommon.P2PVersion032}}, p2pcommon.P2PVersion200},
		{"TMulti", args{[]p2pcommon.P2PVersion{p2pcommon.P2PVersion230, p2pcommon.P2PVersion200}}, p2pcommon.P2PVersion230},
		{"TTooOld", args{[]p2pcommon.P2PVersion{p2pcommon.P2PVersion030}}, p2pcommon.P2PVersionUnknown},
		{"TUnknown", args{[]p2pcommon.P2PVersion{9999999, 9999998}}, p2pcommon.P2PVersionUnknown},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			is := p2pmock.NewMockInternalService(ctrl)
			pm := p2pmock.NewMockPeerManager(ctrl)
			actor := p2pmock.NewMockActorService(ctrl)
			ca := p2pmock.NewMockChainAccessor(ctrl)
			vm := newDefaultVersionManager(is, actor, pm, ca, logger, dummyChainID)

			if got := vm.FindBestP2PVersion(tt.args.versions); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("defaultVersionManager.FindBestP2PVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultVersionManager_GetVersionedHandshaker(t *testing.T) {
	dummyChainID := &types.ChainID{}
	if chain.Genesis == nil {
		chain.Genesis = &types.Genesis{ID: *dummyChainID}
	}

	// must return valid handshaker for versions in global var p2pcommon.AcceptedInboundVersions and p2pcommon.AttemptingOutboundVersions

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	is := p2pmock.NewMockInternalService(ctrl)
	pm := p2pmock.NewMockPeerManager(ctrl)
	actor := p2pmock.NewMockActorService(ctrl)
	ca := p2pmock.NewMockChainAccessor(ctrl)
	r := p2pmock.NewMockReadWriteCloser(ctrl)
	sampleID := types.RandomPeerID()

	ca.EXPECT().ChainID(gomock.Any()).Return(dummyChainID).MaxTimes(1)
	is.EXPECT().CertificateManager().Return(nil).AnyTimes()
	is.EXPECT().SelfMeta().Return(sampleMeta).AnyTimes()

	h := newDefaultVersionManager(is, actor, pm, ca, logger, dummyChainID)

	arrs := append([]p2pcommon.P2PVersion{}, p2pcommon.AcceptedInboundVersions...)
	arrs = append(arrs, p2pcommon.AttemptingOutboundVersions...)

	for _, v := range arrs {
		got, err := h.GetVersionedHandshaker(v, sampleID, r)
		if err != nil {
			t.Fatalf("defaultVersionManager.GetVersionedHandshaker(%v) error = %v, want no error",v,err)
		}
		if got == nil {
			t.Fatalf("defaultVersionManager.GetVersionedHandshaker(%v) return nil, want valid pointer",v)
		}
	}
}

func Test_defaultVersionManager_GetVersionedHandshaker2(t *testing.T) {
	dummyChainID := &types.ChainID{}
	if chain.Genesis == nil {
		chain.Genesis = &types.Genesis{ID:*dummyChainID}
	}


	type args struct {
		version p2pcommon.P2PVersion
	}
	tests := []struct {
		name   string
		args   args

		wantErr bool
	}{
		//
		{"TLegacy", args{p2pcommon.P2PVersion200}, false},
		{"TOld", args{p2pcommon.P2PVersion033}, true},
		{"TOlder", args{p2pcommon.P2PVersion032}, true},
		{"TOldest", args{p2pcommon.P2PVersion031}, true},
		{"TUnknown", args{9999999}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			is := p2pmock.NewMockInternalService(ctrl)
			pm := p2pmock.NewMockPeerManager(ctrl)
			actor := p2pmock.NewMockActorService(ctrl)
			ca := p2pmock.NewMockChainAccessor(ctrl)
			r := p2pmock.NewMockReadWriteCloser(ctrl)
			sampleID := types.RandomPeerID()

			ca.EXPECT().ChainID(gomock.Any()).Return(dummyChainID).MaxTimes(1)
			is.EXPECT().CertificateManager().Return(nil).AnyTimes()
			is.EXPECT().SelfMeta().Return(sampleMeta).AnyTimes()

			h := newDefaultVersionManager(is, actor, pm, ca, logger, dummyChainID)

			got, err := h.GetVersionedHandshaker(tt.args.version, sampleID, r)
			if (err != nil) != tt.wantErr {
				t.Errorf("defaultVersionManager.GetVersionedHandshaker() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got != nil) == tt.wantErr {
				t.Errorf("defaultVersionManager.GetVersionedHandshaker() returns nil == %v , want %v", (got != nil), tt.wantErr )
				return
			}
		})
	}
}
