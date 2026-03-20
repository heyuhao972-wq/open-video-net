package node

import (
	"context"
	"errors"
	"fmt"

	"p2p-node/protocol"
)

func (n *Node) FetchChunk(ctx context.Context, hash string) ([]byte, error) {
	if hash == "" {
		return nil, errors.New("chunk hash required")
	}
	if n.ChunkStore != nil {
		if data, err := n.ChunkStore.Get(hash); err == nil {
			return data, nil
		}
	}

	peers := n.Host.Network().Peers()
	for _, peerID := range peers {
		s, err := n.Host.NewStream(ctx, peerID, protocol.ProtocolID)
		if err != nil {
			continue
		}

		req := protocol.Message{
			Type: "get_chunk",
			Data: []byte(hash),
		}
		if err := protocol.WriteMessage(s, req); err != nil {
			s.Close()
			continue
		}

		resp, err := protocol.ReadMessage(s)
		s.Close()
		if err != nil {
			continue
		}
		if resp.Type != "chunk_data" {
			continue
		}

		if n.ChunkStore != nil {
			_ = n.ChunkStore.Put(hash, resp.Data)
		}
		return resp.Data, nil
	}

	return nil, fmt.Errorf("chunk not found in peers: %s", hash)
}
