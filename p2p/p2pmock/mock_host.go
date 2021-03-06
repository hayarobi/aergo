// Code generated by MockGen. DO NOT EDIT.
// Source: ./vendor/github.com/libp2p/go-libp2p-host/host.go

// Package p2pmock is a generated GoMock package.
package p2pmock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	go_libp2p_interface_connmgr "github.com/libp2p/go-libp2p-interface-connmgr"
	go_libp2p_net "github.com/libp2p/go-libp2p-net"
	go_libp2p_peer "github.com/libp2p/go-libp2p-peer"
	go_libp2p_peerstore "github.com/libp2p/go-libp2p-peerstore"
	go_libp2p_protocol "github.com/libp2p/go-libp2p-protocol"
	go_multiaddr "github.com/multiformats/go-multiaddr"
	go_multistream "github.com/multiformats/go-multistream"
	reflect "reflect"
)

// MockHost is a mock of Host interface
type MockHost struct {
	ctrl     *gomock.Controller
	recorder *MockHostMockRecorder
}

// MockHostMockRecorder is the mock recorder for MockHost
type MockHostMockRecorder struct {
	mock *MockHost
}

// NewMockHost creates a new mock instance
func NewMockHost(ctrl *gomock.Controller) *MockHost {
	mock := &MockHost{ctrl: ctrl}
	mock.recorder = &MockHostMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHost) EXPECT() *MockHostMockRecorder {
	return m.recorder
}

// ID mocks base method
func (m *MockHost) ID() go_libp2p_peer.ID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(go_libp2p_peer.ID)
	return ret0
}

// ID indicates an expected call of ID
func (mr *MockHostMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockHost)(nil).ID))
}

// Peerstore mocks base method
func (m *MockHost) Peerstore() go_libp2p_peerstore.Peerstore {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Peerstore")
	ret0, _ := ret[0].(go_libp2p_peerstore.Peerstore)
	return ret0
}

// Peerstore indicates an expected call of Peerstore
func (mr *MockHostMockRecorder) Peerstore() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Peerstore", reflect.TypeOf((*MockHost)(nil).Peerstore))
}

// Addrs mocks base method
func (m *MockHost) Addrs() []go_multiaddr.Multiaddr {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Addrs")
	ret0, _ := ret[0].([]go_multiaddr.Multiaddr)
	return ret0
}

// Addrs indicates an expected call of Addrs
func (mr *MockHostMockRecorder) Addrs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Addrs", reflect.TypeOf((*MockHost)(nil).Addrs))
}

// Network mocks base method
func (m *MockHost) Network() go_libp2p_net.Network {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Network")
	ret0, _ := ret[0].(go_libp2p_net.Network)
	return ret0
}

// Network indicates an expected call of Network
func (mr *MockHostMockRecorder) Network() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Network", reflect.TypeOf((*MockHost)(nil).Network))
}

// Mux mocks base method
func (m *MockHost) Mux() *go_multistream.MultistreamMuxer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Mux")
	ret0, _ := ret[0].(*go_multistream.MultistreamMuxer)
	return ret0
}

// Mux indicates an expected call of Mux
func (mr *MockHostMockRecorder) Mux() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Mux", reflect.TypeOf((*MockHost)(nil).Mux))
}

// Connect mocks base method
func (m *MockHost) Connect(ctx context.Context, pi go_libp2p_peerstore.PeerInfo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect", ctx, pi)
	ret0, _ := ret[0].(error)
	return ret0
}

// Connect indicates an expected call of Connect
func (mr *MockHostMockRecorder) Connect(ctx, pi interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockHost)(nil).Connect), ctx, pi)
}

// SetStreamHandler mocks base method
func (m *MockHost) SetStreamHandler(pid go_libp2p_protocol.ID, handler go_libp2p_net.StreamHandler) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetStreamHandler", pid, handler)
}

// SetStreamHandler indicates an expected call of SetStreamHandler
func (mr *MockHostMockRecorder) SetStreamHandler(pid, handler interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStreamHandler", reflect.TypeOf((*MockHost)(nil).SetStreamHandler), pid, handler)
}

// SetStreamHandlerMatch mocks base method
func (m *MockHost) SetStreamHandlerMatch(arg0 go_libp2p_protocol.ID, arg1 func(string) bool, arg2 go_libp2p_net.StreamHandler) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetStreamHandlerMatch", arg0, arg1, arg2)
}

// SetStreamHandlerMatch indicates an expected call of SetStreamHandlerMatch
func (mr *MockHostMockRecorder) SetStreamHandlerMatch(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStreamHandlerMatch", reflect.TypeOf((*MockHost)(nil).SetStreamHandlerMatch), arg0, arg1, arg2)
}

// RemoveStreamHandler mocks base method
func (m *MockHost) RemoveStreamHandler(pid go_libp2p_protocol.ID) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveStreamHandler", pid)
}

// RemoveStreamHandler indicates an expected call of RemoveStreamHandler
func (mr *MockHostMockRecorder) RemoveStreamHandler(pid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveStreamHandler", reflect.TypeOf((*MockHost)(nil).RemoveStreamHandler), pid)
}

// NewStream mocks base method
func (m *MockHost) NewStream(ctx context.Context, p go_libp2p_peer.ID, pids ...go_libp2p_protocol.ID) (go_libp2p_net.Stream, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, p}
	for _, a := range pids {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "NewStream", varargs...)
	ret0, _ := ret[0].(go_libp2p_net.Stream)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewStream indicates an expected call of NewStream
func (mr *MockHostMockRecorder) NewStream(ctx, p interface{}, pids ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, p}, pids...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewStream", reflect.TypeOf((*MockHost)(nil).NewStream), varargs...)
}

// Close mocks base method
func (m *MockHost) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockHostMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockHost)(nil).Close))
}

// ConnManager mocks base method
func (m *MockHost) ConnManager() go_libp2p_interface_connmgr.ConnManager {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConnManager")
	ret0, _ := ret[0].(go_libp2p_interface_connmgr.ConnManager)
	return ret0
}

// ConnManager indicates an expected call of ConnManager
func (mr *MockHostMockRecorder) ConnManager() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnManager", reflect.TypeOf((*MockHost)(nil).ConnManager))
}
