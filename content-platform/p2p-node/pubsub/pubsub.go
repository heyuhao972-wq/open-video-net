package pubsub

import (
	"context"

	gosub "github.com/libp2p/go-libp2p-pubsub"
)

func Publish(ps *gosub.PubSub, topic string, msg []byte) error {
	t, err := ps.Join(topic)
	if err != nil {
		return err
	}

	return t.Publish(context.Background(), msg)
}
