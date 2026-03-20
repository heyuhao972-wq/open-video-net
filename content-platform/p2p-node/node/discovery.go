package node

import (
	"context"
	"log"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
)

func SetupDHT(ctx context.Context, h host.Host) (*dht.IpfsDHT, error) {

	kad, err := dht.New(ctx, h)
	if err != nil {
		return nil, err
	}

	err = kad.Bootstrap(ctx)
	if err != nil {
		return nil, err
	}

	log.Println("DHT bootstrapped")

	return kad, nil
}
