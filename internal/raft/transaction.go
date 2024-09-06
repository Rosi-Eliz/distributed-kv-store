package raft

import "time"

// Command represents a single operation in the transaction
type Command struct {
	Op    string        `json:"op"`
	Key   string        `json:"key"`
	Value string        `json:"value"`
	TTL   time.Duration `json:"ttl"`
}

// Transaction represents a group of commands that should be applied together
type Transaction struct {
	Commands []Command `json:"commands"`
}
