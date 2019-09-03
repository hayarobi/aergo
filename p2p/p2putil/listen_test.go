/*
 * @file
 * @copyright defined in aergo/LICENSE.txt
 */

package p2putil

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aergoio/aergo/types"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	nat "github.com/libp2p/go-libp2p-nat"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/libp2p/go-libp2p-peerstore/pstoremem"
	"github.com/multiformats/go-multiaddr"
	"testing"
	"time"
)

func TestListen(t *testing.T) {
	//t.Skip("just for experiment and very enviroment specific, so not suitable for unittest")
	lifeTime := time.Second * 30
	ctx := context.Background()

	samplePKBase64 := "CAISIM1yE7XjJyKTw4fQYMROnlxmEBut5OPPGVde7PeVAf0x"
	bytes, _ := base64.StdEncoding.DecodeString(samplePKBase64)
	samplePrivKey, _ := crypto.UnmarshalPrivateKey(bytes)
	fmt.Println("Staring tests... ")

	listens := make([]types.Multiaddr, 0, 2)
	listen, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/7847")
	//listen, err := multiaddr.NewMultiaddr("/dns4/iparkmac.blocko.io/tcp/7847")
	if err != nil {
		t.Fatalf("Error while make multiaddr: %v", err)
	}
	listens = append(listens, listen)
	addressFactory := func(addrs []types.Multiaddr) []types.Multiaddr {
		addrs = append(addrs, listen)
		return addrs
	}
	if false {
		addressFactory(nil)
	}

	fmt.Println("Running listen host... ")

	peerStore := pstore.NewPeerstore(pstoremem.NewKeyBook(), pstoremem.NewAddrBook(), pstoremem.NewProtoBook(), pstoremem.NewPeerMetadata())
	opts := []libp2p.Option{
		libp2p.Identity(samplePrivKey), libp2p.Peerstore(peerStore), libp2p.ListenAddrs(listens...),
		//libp2p.AddrsFactory(addressFactory),
		libp2p.NATPortMap(),
	}
	newHost, err := libp2p.New(ctx, opts...)
	if err != nil {
		t.Errorf("Error while start listening: %v", err)
	}
	newHost.SetStreamHandler("/", func(stream network.Stream) {
		fmt.Printf("New conn from %v ", stream.Conn().RemoteMultiaddr().String())
	})

	fmt.Printf("Addrs are %v \n", newHost.Addrs())
	nat, err := nat.DiscoverNAT(ctx)
	if err != nil {
		fmt.Printf("Failed to discover nat\n")
	}
	fmt.Println("NAT mappings... ")
	if nat == nil {
		fmt.Println("NAT is nil")
	} else {
		aergoMapping, err := nat.NewMapping("tcp", 7847)
		if err != nil {
			fmt.Printf("Failed to add nat mapping %v\n",err)
		} else {
			fmt.Printf("Added mapping %v \n", aergoMapping)
		}
		for i, m := range nat.Mappings() {
			fmt.Printf("m%v, %v\n",i, m)
		}
	}
	fmt.Println("Sleeping for a while")
	time.Sleep(lifeTime)
	newHost.Close()
	fmt.Printf("Finished!")
}
