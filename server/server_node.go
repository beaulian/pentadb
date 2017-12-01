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

package main

import (
	"sync"
	"errors"
	"math/rand"
	"fmt"

	"github.com/shenaishiren/pentadb/args"
	"github.com/shenaishiren/pentadb/log"
	"github.com/syndtr/goleveldb/leveldb"
)

var LOG = log.DefaultLog

type ServerNode struct {
	Ipaddr string

	OtherNodes []string

	ReplicaNodes []string

	DB *leveldb.DB

	mutex *sync.RWMutex   // read-write lock
}

func NewServerNode(ipaddr string) *ServerNode {
	return &ServerNode {
		Ipaddr: ipaddr,
		mutex: new(sync.RWMutex),
	}
}

func (n *ServerNode) randomChoice(list []string, k int) []string {
	pool := list
	p := len(pool)
	result := make([]string, k)
	for i := 0; i < k; i++ {
		j := rand.Intn(p - i)
		result[i] = pool[j]
		pool[j] = result[p - i - 1]
	}
	return result
}

func (n *ServerNode) Init(args *args.InitArgs, result *[]byte) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	n.Ipaddr = args.Self
	n.OtherNodes = args.OtherNodes
	replicaNodes := n.randomChoice(args.OtherNodes, args.Replicas)
	if len(replicaNodes) == 0 {
		return errors.New(fmt.Sprintf("node %s init failed", n.Ipaddr))
	}
	replicaNodes = append(replicaNodes, n.Ipaddr)
	n.ReplicaNodes = replicaNodes
	return nil
}

func (n *ServerNode) AddNode(node string, result *[]byte) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	n.OtherNodes = append(n.OtherNodes, node)
	return nil
}

func (n *ServerNode) RemoveNode(node string, result *[]byte) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	var j = 0
	var flag = false
	for _, otherNode := range n.OtherNodes {
		if otherNode != node {
			n.OtherNodes[j] = otherNode
			j++
		} else {
			flag = true
		}
	}
	if flag {
		n.OtherNodes = n.OtherNodes[0:len(n.OtherNodes) - 1]
	}
	return nil
}

func (n *ServerNode) Put(args *args.KVArgs, result *[]byte) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	err := n.DB.Put(args.Key, args.Value, nil)
	return err
}

func (n *ServerNode) Get(key []byte, result *[]byte) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	res, err := n.DB.Get(key, nil)
	*result = res
	return err
}

func (n *ServerNode) Delete(key []byte, result *[]byte) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	err := n.DB.Delete(key, nil)
	return err
}