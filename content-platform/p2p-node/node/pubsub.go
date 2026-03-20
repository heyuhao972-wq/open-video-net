package node

import (
	"log"

	ps "p2p-node/pubsub"
)

func (n *Node) Publish(topic string, msg []byte) {
	if err := ps.Publish(n.PubSub, topic, msg); err != nil {
		log.Println("publish error:", err)
		return
	}
}
