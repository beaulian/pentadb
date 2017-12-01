// Contains the interface and implementation of PentaDB Client

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
	"fmt"
	"errors"
	"github.com/shenaishiren/pentadb/log"
)

const (
	defaultReplicas = 0                 // default replicas for raft algorithm
	maxn = 1024
)

var LOG = log.DefaultLog

type Client struct {
	// all nodes in hash ring
	nodes map[string]*ClientNode

	// hash ring
	hashRing *HashRing

	// the channel is used for checking what node is lost
	unreachableChan chan string
}

func NewClient(nodeIpaddrs []string, weights map[string]int, replicas int) (*Client, error) {
	// check nodes' count
	nodesCount := len(nodeIpaddrs)
	// TODO
	if nodesCount < defaultReplicas + 1 {
		return nil, errors.New(
			fmt.Sprintf("nodes must be > %d", defaultReplicas),
		)
	}
	// TODO
	if replicas > nodesCount || replicas < defaultReplicas {
		return nil, errors.New(
			fmt.Sprintf("replicas must > %d and < %d", defaultReplicas, nodesCount),
		)
	}
	// initialize hash ring
	hashRing := NewHashRing()
	// check whether the ip address is connectable or not in `init` function
	nodes, _ := hashRing.init(nodeIpaddrs, weights)
	nodeDict := make(map[string]*ClientNode)
	// initialize client
	client := &Client{
		nodes: nodeDict,
		hashRing: hashRing,
		unreachableChan: make(chan string, maxn),
	}
	for _, node := range nodes {
		nodeDict[node.Name] = node
		// asynchronously
		go node.Proxy.Init(nodeIpaddrs, replicas, client.unreachableChan)
	}
	// event loop about checking nodes
	go func() {
		for {
			select {
				case nodeName := <- client.unreachableChan:
					client.RemoveNode(nodeName)
			}
		}
	}()

	return client, nil
}

func (c *Client) AddNode(nodeIpaddr string, weight int) {
	node := c.hashRing.addNode(nodeIpaddr, weight)
	if node == nil {
		return
	}
	c.nodes[nodeIpaddr] = node
	node.Proxy.AddNode(nodeIpaddr, c.unreachableChan)
}

func (c *Client) RemoveNode(nodeName string) {
	node := c.nodes[nodeName]
	if node != nil {
		c.hashRing.deleteNode(node.Ipaddr, node.Weight)
		go node.Proxy.RemoveNode(node.Ipaddr, c.unreachableChan)
		delete(c.nodes, nodeName)
	}
}

func (c *Client) Put(key []byte, value []byte) {
	// choose a node
	hashKey := KemataHash(Md5Hash(key), 0)
	node, err := c.hashRing.findProperNode(hashKey)
	if err != nil {
		LOG.Error("error occurred when find proper node: ", err.Error())
		return
	}
	node.rNode.Proxy.Put(key, value, c.unreachableChan)
}

func (c *Client) Get(key []byte) []byte {
	hashKey := KemataHash(Md5Hash(key), 0)
	node, err := c.hashRing.findProperNode(hashKey)
	if err != nil {
		LOG.Error("error occurred when find proper node: ", err.Error())
		return nil
	}
	return node.rNode.Proxy.Get(key, c.unreachableChan)
}

func (c *Client) Delete(key []byte) {
	hashKey := KemataHash(Md5Hash(key), 0)
	node, err := c.hashRing.findProperNode(hashKey)
	if err != nil {
		LOG.Error("error occurred when find proper node: ", err.Error())
		return
	}
	node.rNode.Proxy.Delete(key, c.unreachableChan)
}

func (c *Client) Close() {
	close(c.unreachableChan)
}