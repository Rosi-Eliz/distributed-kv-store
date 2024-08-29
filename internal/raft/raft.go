package raft

import (
	"github.com/Rosi-Eliz/distributed-kv-store/internal/store"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"os"
	"path/filepath"
	"time"
)

type Node struct {
	Raft *raft.Raft
	FSM  *FSM
}

func NewRaftNode(nodeID string, bindAddr string, raftDir string, kvStore *store.Store) (*Node, error) {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeID)

	transport, err := raft.NewTCPTransport(bindAddr, nil, 3, 10*time.Second, os.Stdout)
	if err != nil {
		return nil, err
	}

	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(raftDir, "raft-stable.db"))
	if err != nil {
		return nil, err
	}

	logStore, err := raftboltdb.NewBoltStore(filepath.Join(raftDir, "raft-log.db"))
	if err != nil {
		return nil, err
	}

	snapshotStore, err := raft.NewFileSnapshotStore(raftDir, 1, os.Stdout)
	if err != nil {
		return nil, err
	}

	fsm := NewFSM(kvStore)

	ra, err := raft.NewRaft(config, fsm, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		return nil, err
	}

	return &Node{Raft: ra}, nil
}
