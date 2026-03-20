package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"p2p-node/config"
	"p2p-node/node"
)

func main() {

	ctx := context.Background()

	cfg := config.LoadRuntimeConfig()
	n, err := node.NewNode(ctx, cfg.ChunkDir)
	if err != nil {
		panic(err)
	}

	fmt.Println("Node started")
	fmt.Println("PeerID:", n.Host.ID())

	for _, addr := range n.Host.Addrs() {
		fmt.Println("Address:", addr)
	}

	go n.StartDiscovery(ctx)
	go startHTTPServer(ctx, cfg.HTTPPort, n)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		cmd, _ := reader.ReadString('\n')
		cmd = strings.TrimSpace(cmd)
		if cmd == "" {
			continue
		}

		args := strings.Split(cmd, " ")
		if len(args) == 0 || args[0] == "" {
			continue
		}

		switch args[0] {

		case "peers":

			peers := n.Host.Network().Peers()

			for _, p := range peers {
				fmt.Println(p)
			}

		case "msg":

			if len(args) < 2 {
				continue
			}

			msg := strings.Join(args[1:], " ")

			n.Publish("chat", []byte(msg))

		}
	}

}

func startHTTPServer(ctx context.Context, port string, n *node.Node) {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/chunk/", func(w http.ResponseWriter, r *http.Request) {
		hash := strings.TrimPrefix(r.URL.Path, "/chunk/")
		if hash == "" {
			http.Error(w, "missing hash", http.StatusBadRequest)
			return
		}
		data, err := n.FetchChunk(ctx, hash)
		if err != nil {
			http.Error(w, "chunk not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	})

	addr := ":" + port
	fmt.Println("P2P HTTP gateway running on port:", port)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Println("http server error:", err)
	}
}
