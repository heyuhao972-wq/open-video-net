package node

import (
	"context"
	"log"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/multiformats/go-multiaddr"
)

var bootstrapPeers = []string{
	"/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
}

func (n *Node) ConnectBootstrap(ctx context.Context) error {

	for _, addr := range bootstrapPeers {

		maddr, err := multiaddr.NewMultiaddr(addr)

		if err != nil {
			continue
		}

		info, err := peer.AddrInfoFromP2pAddr(maddr)

		if err != nil {
			continue
		}

		if err := n.Host.Connect(ctx, *info); err != nil {
			log.Println("bootstrap connect error:", addr, err)
		}
	}

	return nil
}
