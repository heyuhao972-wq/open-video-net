package node

import "fmt"

func (n *Node) ListPeers() {

	peers := n.Host.Network().Peers()

	for _, p := range peers {

		fmt.Println(p)

	}
}
