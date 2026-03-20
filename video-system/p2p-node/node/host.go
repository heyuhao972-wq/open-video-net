package node

import (
	"context"

	libp2p "github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
)

func CreateHost(ctx context.Context) (host.Host, error) {

	return libp2p.New()
}
