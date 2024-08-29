package raft

import (
	"github.com/hashicorp/raft"
	"log"
)

func (r *Node) AddVoter(id, address string) error {
	future := r.Raft.AddVoter(raft.ServerID(id), raft.ServerAddress(address), 0, 0)
	if future.Error() != nil {
		return future.Error()
	}
	log.Printf("node %s at %s added as voter", id, address)
	return nil
}

func (r *Node) RemoveVoter(id string) error {
	future := r.Raft.RemoveServer(raft.ServerID(id), 0, 0)
	if future.Error() != nil {
		return future.Error()
	}
	log.Printf("node %s removed as voter", id)
	return nil
}
