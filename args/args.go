package args

// common variable
type InitArgs struct {
	Self string

	// a array consisting of other nodes' ipaddr
	OtherNodes []string

	// replicas
	Replicas int
}

type KVArgs struct {
	Key []byte

	Value []byte
}
