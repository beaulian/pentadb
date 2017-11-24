// Contains the interface and implementation of Node
// The `Node` struct is designed for client which is distinguished
// with the `Node` struct in server. When a node is added to cluster,
// client will maintain a connection with every node, so it's necessary
// to record the status of each node.

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
	"time"

	"github.com/satori/go.uuid"
)

var timeout = 2 * time.Second


type Node struct {
	// The name of node, maybe used for some commands
	// like `pentadb nodes list`
	Name string

	// The address of node, format as `<ip>:<port>`
	Ipaddr string

	// creating time
	Ctime time.Time

	// Node Proxy
	Proxy *NodeProxy
}

func NewNode(ipaddr string) *Node {
	node := &Node {
		Name:     uuid.NewV1().String(),
		Ipaddr:   ipaddr,
		Ctime:    time.Now(),
	}
	if !Reachable(ipaddr, timeout) {
		return nil
	}
	proxy := NewNodeProxy(node)
	if proxy == nil {
		return nil
	}
	node.Proxy = proxy
	return node
}