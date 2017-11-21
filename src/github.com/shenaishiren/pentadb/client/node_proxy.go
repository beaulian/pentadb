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

package client

import (
	"fmt"
	"errors"
	"net/rpc"

	"github.com/shenaishiren/pentadb/opt"
	"github.com/shenaishiren/pentadb/args"
)

type NodeProxy struct {
	// client-side node
	node *Node

	// rpc client
	rpcClient *rpc.Client

	// node's state
	nodeState opt.NodeState
}

func NewNodeProxy(node *Node) *NodeProxy {
	client, err := rpc.DialHTTP(opt.DefaultProtocol, node.Ipaddr)
	if err != nil {
		return nil
	}
	return &NodeProxy{
		node:          node,
		nodeState:     opt.NodeRunning,
		rpcClient:     client,
	}
}

func (np *NodeProxy) init(nodeIpaddrs []string, option *opt.NodeProxyOptions) error {
	nodesCount := len(nodeIpaddrs)
	replicas := option.Replicas
	if option == nil || option.Replicas == 0 {
		replicas = opt.DefaultReplicas
	} else if replicas < opt.DefaultReplicas || replicas >= nodesCount {
		return errors.New(
			fmt.Sprintf("replica number must be >= %d and < %d",
				opt.DefaultReplicas, nodesCount),
		)
	}
	args := &args.InitArgs{
		Nodes: nodeIpaddrs, Replicas: replicas,
	}

	return np.rpcClient.Call("Node.init", args, nil)
}

func (np *NodeProxy) addNode(nodeIpaddr string) error {
	return np.rpcClient.Call("Node.addNode", nodeIpaddr, nil)
}

func (np *NodeProxy) removeNode(nodeIpaddr string) error {
	return np.rpcClient.Call("Node.removeNode", nodeIpaddr, nil)
}

func (np *NodeProxy) put(key []byte, value []byte) error {
	kvArgs := &args.KVArgs{key, value}
	return np.rpcClient.Call("Node.put", kvArgs, nil)
}

func (np *NodeProxy) get(key []byte) ([]byte, error) {
	var result []byte
	err := np.rpcClient.Call("Node.get", key, &result)
	return result, err
}

func (np *NodeProxy) delete(key []byte) error {
	return np.rpcClient.Call("Node.delete", key, nil)
}

//func (np *NodeProxy) setReplicas(rNode *Node, replica int) error {
//	// we adopt a random method to choose N-1 nodes
//	// which makes it fairer, then append itself to the
//	// N-1 nodes
//	replicas := make([]*Node, replica)
//	// rNode is a closure
//	filterFunc := func(i interface{}) bool {
//		return i.(*Node).Ipaddr != rNode.Ipaddr
//	}
//	temp, _ := utils.RandomChoice(utils.Filter(filterFunc, np.hr.nodes.Values()), replica)
//	for _, v := range temp {
//		replicas = append(replicas, v.(*Node))
//	}
//	replicas = append(replicas, rNode)
//	rNode.Replicas = replicas
//
//	return nil
//}
//
//func (np *NodeProxy) getProperNode(key string) (*Node, error) {
//	hashKey := KemataHash(key, 0)
//	vNode, err := np.hr.findProperNode(hashKey)
//
//	return vNode.RNode, err
//}
