package node

import (
	"context"
	"log"

	libp2p "github.com/libp2p/go-libp2p"

	pubsub "github.com/libp2p/go-libp2p-pubsub"

	dht "github.com/libp2p/go-libp2p-kad-dht"

	"github.com/libp2p/go-libp2p/core/host"

	discovery "github.com/libp2p/go-libp2p/p2p/discovery/routing"

	"github.com/libp2p/go-libp2p/p2p/discovery/util"

	"p2p-node/chunk"
	"p2p-node/protocol"
)

type Node struct {
	Host host.Host

	DHT *dht.IpfsDHT

	PubSub *pubsub.PubSub

	Discovery *discovery.RoutingDiscovery

	ChunkStore *chunk.Store
}

func NewNode(ctx context.Context, chunkDir string) (*Node, error) {

	h, err := libp2p.New()
	if err != nil {
		return nil, err
	}

	kad, err := dht.New(ctx, h)
	if err != nil {
		return nil, err
	}

	err = kad.Bootstrap(ctx)
	if err != nil {
		return nil, err
	}

	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		return nil, err
	}

	disc := discovery.NewRoutingDiscovery(kad)

	n := &Node{
		Host:      h,
		DHT:       kad,
		PubSub:    ps,
		Discovery: disc,
	}

	store, err := chunk.NewStore(chunkDir)
	if err != nil {
		return nil, err
	}
	n.ChunkStore = store

	h.SetStreamHandler(protocol.ProtocolID, n.HandleStream)

	err = n.ConnectBootstrap(ctx)

	if err != nil {
		log.Println("bootstrap error:", err)
	}

	return n, nil
}

func (n *Node) StartDiscovery(ctx context.Context) {

	util.Advertise(ctx, n.Discovery, "video-network")

	peerChan, err := n.Discovery.FindPeers(ctx, "video-network")
	if err != nil {
		log.Println("find peers error:", err)
		return
	}

	for peer := range peerChan {

		if peer.ID == n.Host.ID() {
			continue
		}

		if err := n.Host.Connect(ctx, peer); err != nil {
			log.Println("connect discovered peer error:", peer.ID, err)
		}
	}

}
