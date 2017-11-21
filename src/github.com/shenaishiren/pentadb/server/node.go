// Contains the interface and implementation of Node

/* BSD 3-Clause License

Copyright (c) 2017, Guan Jiawen, Li Lundong
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of the copyright holder nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package server

import (
	"sync"
	"errors"
	"math/rand"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/shenaishiren/pentadb/args"
)

type NodeStatus int

const (
	Running NodeStatus = iota
	Terminal
)

type Node struct {
	ipaddr string

	state NodeStatus

	otherNodes []string

	replicaNodes []string

	db *leveldb.DB

	mutex *sync.Mutex
}

func NewNode(ipaddr string) *Node {
	return &Node {
		ipaddr: ipaddr,
		state: Running,
		mutex: new(sync.Mutex),
	}
}

func (n *Node) randomChoice(list []string, k int) ([]string, error) {
	if k <= 0 {
		return nil, errors.New("invalid k: k must be > 0")
	}
	pool := list
	p := len(pool)
	result := make([]string, k)
	for i := 0; i < k; i++ {
		j := rand.Intn(p - i)
		result[i] = pool[j]
		pool[j] = result[p - i - 1]
	}

	return result, nil
}


func (n *Node) init(args *args.InitArgs, result *[]byte) error {
	nodes := args.Nodes
	otherNodes := make([]string, len(nodes) - 1)
	for _, node := range nodes {
		if node != n.ipaddr {
			otherNodes = append(otherNodes, node)
		}
	}

	n.otherNodes = otherNodes
	replicaNodes, err := n.randomChoice(otherNodes, args.Replicas)
	if err != nil {
		return err
	}
	replicaNodes = append(replicaNodes, n.ipaddr)
	n.replicaNodes = replicaNodes
	return nil
}

func (n *Node) addNode(node string, result *[]byte) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	n.otherNodes = append(n.otherNodes, node)
}

func (n *Node) removeNode(node string, result *error) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	var j = 0
	var flag = false
	for _, otherNode := range n.otherNodes {
		if otherNode != node {
			n.otherNodes[j] = otherNode
			j++
		} else {
			flag = true
		}
	}
	if flag {
		n.otherNodes = n.otherNodes[0:len(n.otherNodes) - 1]
	}
}

func (n *Node) put(args *args.KVArgs, result *[]byte) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	err := n.db.Put(args.Key, args.Value, nil)
	return err
}

func (n *Node) get(key []byte, result *[]byte) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	res, err := n.db.Get(key, nil)
	result = &res
	return err
}

func (n *Node) delete(key []byte, result *[]byte) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	err := n.db.Delete(key, nil)
	return err
}