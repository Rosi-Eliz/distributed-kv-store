package main

import (
	"github.com/Rosi-Eliz/distributed-kv-store/internal/raft"
	"github.com/Rosi-Eliz/distributed-kv-store/internal/server"
	"github.com/Rosi-Eliz/distributed-kv-store/internal/store"
	"log"
)

func main() {
	kvStore := store.NewStore()

	raftNode, err := raft.NewRaftNode("node1", "localhost:5001", "./raft", kvStore)
	if err != nil {
		log.Fatalf("failed to create raft node: %v", err)
	}

	srv := server.NewServer(raftNode)
	srv.Run(":8080")
}
