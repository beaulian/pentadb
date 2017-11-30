package opt

import "time"

const (
	DefaultReplicas = 1                 // default replicas for raft algorithm
	DeafultPath = "/tmp/pentadb"        // default path for levelDB
	DefaultProtocol = "tcp"
	DefaultTimeout = 3 * time.Second
)

type NodeState int

const (
	NodeRunning NodeState = iota
	NodeTerminal
)
