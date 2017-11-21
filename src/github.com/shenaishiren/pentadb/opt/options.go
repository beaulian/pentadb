package opt

const (
	DefaultReplicas = 3                 // default replicas for raft algorithm
	DeafultPath = "/usr/local/pentadb"  // default path for levelDB
	DefaultProtocol = "tcp"
)

type NodeState int

const (
	NodeRunning NodeState = iota
	NodeTerminal
)

type NodeProxyOptions struct {
	// node's replicas
	Replicas int

	// Whether to force the flush to the node
	Flush bool
}
