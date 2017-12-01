// Contains the interface and implementation of Node Proxy
// which is the manager of a node

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

package pentadb

import (
	"sync"
	"time"

	"github.com/shenaishiren/pentadb/args"
	nrpc "github.com/shenaishiren/pentadb/rpc"
)

const (
	defaultProtocol = "tcp"
	defaultTimeout = 3 * time.Second
)
// end

type NodeProxy struct {
	// client-side node
	node *ClientNode

	mu *sync.Mutex
}

func NewNodeProxy(node *ClientNode) *NodeProxy {
	if !Reachable(node.Ipaddr, defaultTimeout) {
		return nil
	}
	return &NodeProxy{
		node:          node,
		mu:            new(sync.Mutex),
	}
}

func (np *NodeProxy) call(serviceMethod string, args interface{}, unreachableChan chan string) []byte {
	client, err := nrpc.DialTimeout(defaultProtocol, np.node.Ipaddr, defaultTimeout)
	if err != nil {
		LOG.Errorf("node %s is unreachable: %s", np.node.Ipaddr, err.Error())
		unreachableChan <- np.node.Name
		return nil
	}
	// call
	var result []byte
	err = client.Call(serviceMethod, args, &result)
	if err != nil {
		LOG.Error("rpc call failed: ", err.Error())
		return nil
	}
	if err = client.Close(); err != nil {
		LOG.Error("client close failed: ", err.Error())
	}
	return result
}

func (np *NodeProxy) Init(nodeIpaddrs []string, replicas int, unreachableChan chan string) {
	var otherNodes []string
	for _, node := range nodeIpaddrs {
		if node != np.node.Ipaddr {
			otherNodes = append(otherNodes, node)
		}
	}
	args := &args.InitArgs{
		Self: np.node.Ipaddr,
		OtherNodes: otherNodes,
		Replicas: replicas,
	}
	np.call("ServerNode.Init", args, unreachableChan)
}

func (np *NodeProxy) AddNode(nodeIpaddr string, unreachableChan chan string) {
	np.call("ServerNode.AddNode", nodeIpaddr, unreachableChan)
}

func (np *NodeProxy) RemoveNode(nodeIpaddr string, unreachableChan chan string) {
	np.call("ServerNode.RemoveNode", nodeIpaddr, unreachableChan)
}

func (np *NodeProxy) Put(key []byte, value []byte, unreachableChan chan string) {
	kvArgs := &args.KVArgs{Key:key, Value: value}
	np.call("ServerNode.Put", kvArgs, unreachableChan)
}

func (np *NodeProxy) Get(key []byte, unreachableChan chan string) []byte {
	return np.call("ServerNode.Get", key, unreachableChan)
}

func (np *NodeProxy) Delete(key []byte, unreachableChan chan string) {
	np.call("ServerNode.Delete", key, unreachableChan)
}