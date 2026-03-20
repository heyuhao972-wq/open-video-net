package node

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/network"

	"p2p-node/protocol"
)

func (n *Node) HandleStream(s network.Stream) {
	fmt.Println("stream connected")
	defer s.Close()

	msg, err := protocol.ReadMessage(s)
	if err != nil {
		fmt.Println("read stream error:", err)
		return
	}

	switch msg.Type {
	case "get_chunk":
		hash := string(msg.Data)
		if hash == "" {
			_ = protocol.WriteMessage(s, protocol.Message{
				Type: "error",
				Data: []byte("missing chunk hash"),
			})
			return
		}
		data, err := n.ChunkStore.Get(hash)
		if err != nil {
			_ = protocol.WriteMessage(s, protocol.Message{
				Type: "error",
				Data: []byte("chunk not found"),
			})
			return
		}
		_ = protocol.WriteMessage(s, protocol.Message{
			Type: "chunk_data",
			Data: data,
		})
	default:
		fmt.Println("unknown message type:", msg.Type)
		_ = protocol.WriteMessage(s, protocol.Message{
			Type: "error",
			Data: []byte("unknown message type"),
		})
	}
}
