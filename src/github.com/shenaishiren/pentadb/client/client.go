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

package client

import (
	"fmt"
	"errors"
	"github.com/shenaishiren/pentadb/opt"
	"github.com/shenaishiren/pentadb/log"
)

var LOG = log.DefaultLog

const MAXN = 1024

type Client struct {
	// all nodes in hash ring
	nodes map[string]*Node

	// hash ring
	hashRing *HashRing

	// the channel is used for checking what node is lost
	unreachableChan chan string
}

func NewClient(nodeIpaddrs []string, weights map[string]int, replicas int) (*Client, error) {
	// check nodes' count
	nodesCount := len(nodeIpaddrs)
	// TODO
	if nodesCount < opt.DefaultReplicas + 1 {
		return nil, errors.New(
			fmt.Sprintf("nodes must be > %d", opt.DefaultReplicas),
		)
	}
	// TODO
	if replicas > nodesCount || replicas < opt.DefaultReplicas {
		return nil, errors.New(
			fmt.Sprintf("replicas must > %d and < %d", opt.DefaultReplicas, nodesCount),
		)
	}
	// initialize hash ring
	hashRing := NewHashRing()
	// check whether the ip address is connectable or not in `init` function
	nodes, _ := hashRing.init(nodeIpaddrs, weights)
	nodeDict := make(map[string]*Node)
	for _, node := range nodes {
		nodeDict[node.Name] = node
		// asynchronously
		go node.Proxy.init(nodeIpaddrs, &opt.NodeProxyOptions{Replicas: replicas})
	}
	// initialize client
	client := &Client{
		nodes: nodeDict,
		hashRing: hashRing,
		unreachableChan: make(chan string, MAXN),
	}
	// event loop about checking nodes
	go func() {
		for {
			select {
				case nodeName := <- client.unreachableChan:
					client.removeNode(nodeName)
			}
		}
	}()

	return client, nil
}

func (c *Client) addNode(nodeIpaddr string, weight int) {
	node := c.hashRing.addNode(nodeIpaddr, weight)
	if node == nil {
		return
	}
	c.nodes[nodeIpaddr] = node
	if err := node.Proxy.addNode(nodeIpaddr); err != nil {
		LOG.Warningf("node: %s is unconnectable, due to error: %s", node.Ipaddr, err.Error())
		c.unreachableChan <- node.Name
	}
}

func (c *Client) removeNode(nodeName string) {
	node := c.nodes[nodeName]
	if err := node.Proxy.removeNode(node.Ipaddr); err != nil {
		LOG.Warningf("node: %s is unconnectable, due to error: %s", node.Ipaddr, err.Error())
	}
	delete(c.nodes, nodeName)
}

func (c *Client) getNode(nodeName string) *Node {
	return c.nodes[nodeName]
}

func (c *Client) put(key []byte, value []byte) error {
	// choose a node
	hashKey := KemataHash(Md5Hash(key), 0)
	node, err := c.hashRing.findProperNode(hashKey)
	if err != nil {
		LOG.Error("error occurred when find proper node: ", err.Error())
		return err
	}
	if err := node.rNode.Proxy.put(key, value); err != nil {
		LOG.Warningf("node: %s is unconnectable, due to error: %s", node.rNode.Ipaddr, err.Error())
		c.unreachableChan <- node.rNode.Name
		return err
	}
	return nil
}

func (c *Client) get(key []byte) ([]byte, error) {
	hashKey := KemataHash(Md5Hash(key), 0)
	node, _ := c.hashRing.findProperNode(hashKey)
	value, err := node.rNode.Proxy.get(key)
	if err != nil {
		LOG.Warningf("node: %s is unconnectable, due to error: %s", node.rNode.Ipaddr, err.Error())
		c.unreachableChan <- node.rNode.Name
		return nil, err
	}
	return value, nil
}

func (c *Client) delete(key []byte) error {
	hashKey := KemataHash(Md5Hash(key), 0)
	node, _ := c.hashRing.findProperNode(hashKey)
	if err := node.rNode.Proxy.delete(key); err != nil {
		LOG.Warningf("node: %s is unconnectable, due to error: %s", node.rNode.Ipaddr, err.Error())
		c.unreachableChan <- node.rNode.Name
		return err
	}
	return nil
}