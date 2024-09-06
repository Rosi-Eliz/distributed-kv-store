package raft

import (
	"bytes"
	"encoding/json"
	"github.com/Rosi-Eliz/distributed-kv-store/internal/store"
	"github.com/hashicorp/raft"
	"io"
	"log"
)

type FSM struct {
	store *store.Store
}

func NewFSM(store *store.Store) *FSM {
	return &FSM{
		store: store,
	}
}

func (f *FSM) Apply(logEntry *raft.Log) interface{} {
	var tx Transaction
	if err := json.NewDecoder(bytes.NewBuffer(logEntry.Data)).Decode(&tx); err != nil {
		log.Println("error decoding transaction")
		return nil
	}

	// Apply all commands in the transaction atomically
	for _, cmd := range tx.Commands {
		switch cmd.Op {
		case "set":
			f.store.Set(cmd.Key, cmd.Value, cmd.TTL)
		case "delete":
			f.store.Delete(cmd.Key)
		}
	}

	return nil
}

func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return &fsmSnapshot{store: f.store}, nil
}

// Restore is crucial for bringing a node back up to date after it reboots or rejoins the cluster.
func (f *FSM) Restore(rc io.ReadCloser) error {
	var data map[string]string
	if err := json.NewDecoder(rc).Decode(&data); err != nil {
		return err
	}

	f.store.ReplaceData(data)

	return nil
}

func (f *FSM) Get(key string) (string, bool) {
	return f.store.Get(key)
}

type fsmSnapshot struct {
	store *store.Store
}

func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		data, err := f.store.SerializeData()
		if err != nil {
			return err
		}

		if _, err := sink.Write(data); err != nil {
			return err
		}

		return nil
	}()

	if err != nil {
		err := sink.Cancel()
		if err != nil {
			return err
		}
		return err
	}
	return sink.Close()
}

func (f *fsmSnapshot) Release() {}
