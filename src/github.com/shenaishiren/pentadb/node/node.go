// Contains the interface and implementation of Node

package node

import (
	"github.com/syndtr/goleveldb/leveldb"
)

type NodeStatus int

const (
	Running NodeStatus = iota
	Terminal
)

const (
	Path = "/var/pentadb"
	defaultReplicas = 3
)

type Node struct {
	Name string         // name
	State NodeStatus   // status
	Replicas []*Node     // replicas
	Hash uint32         // identifier
	DB *leveldb.DB      // database

	Forward []*Node          // point to forward nodes
}

func NewNode(name string, hash uint32, level int) *Node {
	return &Node{
		Name: name,
		State: Running,
		Hash: hash,
		Replicas: make([]*Node, defaultReplicas),
		Forward: make([]*Node, level),
	}
}

