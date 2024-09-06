package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/Rosi-Eliz/distributed-kv-store/internal/raft"
	"github.com/Rosi-Eliz/distributed-kv-store/internal/server"
	"github.com/Rosi-Eliz/distributed-kv-store/internal/store"
	"log"
	"net/http"
	"time"
)

func main() {
	nodeID := flag.String("id", "node1", "Node ID")
	bindAddr := flag.String("addr", ":5001", "Bind address")
	raftDir := flag.String("dir", "./raft1", "Raft storage directory")
	joinAddr := flag.String("join", "", "Join address of an existing node")

	flag.Parse()

	kvStore := store.NewStore()

	raftNode, err := raft.NewRaftNode(*nodeID, *bindAddr, *raftDir, kvStore)
	if err != nil {
		log.Fatalf("failed to create raft node: %v", err)
	}

	srv := server.NewServer(raftNode)
	go srv.Run(*bindAddr)

	// If a join address is provided, join the cluster by making a POST request to /join
	if *joinAddr != "" {
		time.Sleep(2 * time.Second)

		joinReq := struct {
			ID      string `json:"id"`
			Address string `json:"address"`
		}{
			ID:      *nodeID,
			Address: *bindAddr,
		}
		reqBody, err := json.Marshal(joinReq)
		if err != nil {
			log.Fatalf("failed to marshal join request: %v", err)
		}

		// Make the HTTP POST request to the /join endpoint of the existing node
		resp, err := http.Post("http://"+*joinAddr+"/join", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			log.Fatalf("failed to join cluster: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("failed to join cluster, status code: %d", resp.StatusCode)
		}

		log.Printf("Successfully joined the cluster via %s", *joinAddr)
	}

	select {}
}
