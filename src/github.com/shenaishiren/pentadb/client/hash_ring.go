// Contains the interface and implementation of Hash Ring
// implemented by skip list

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
	"math/rand"
	"math"
	"errors"
	"bytes"
	"strings"
	"github.com/seiflotfy/cuckoofilter"
)

const (
	defaultFactor = 40
	defaultWeight = 1
	maxLevel = 32
)

type VNode struct {
	Hash uint32          // identifier
	Forward []*VNode     // point to forward nodes
	rNode *Node          // real server, in order to contact the physical server
}


type HashRing struct {
	rnd *rand.Rand

	header *VNode                     // point to server skip-list
	level int                         // the level of skip-list
	length int                        // number of nodes
	averageWeight float64             // total weight of nodes in hash ring
	filter *cuckoofilter.CuckooFilter // cuckoo filter, ensure every node is unique
}

func NewVNode(node *Node, hash uint32, level int) *VNode {
	return &VNode {
		Hash:    hash,
		Forward: make([]*VNode, level),
		rNode:   node,
	}
}

func NewHashRing() *HashRing {
	return &HashRing{
		rnd:            rand.New(rand.NewSource(0xdeadbeef)),
		level:          1,
		length:         0,
		averageWeight:  0,
		header:         NewVNode(nil,0, maxLevel),
		filter:         cuckoofilter.NewCuckooFilter(100),
	}
}

// Get the count of virtual nodes
func (hr *HashRing) getVNodeCount(weight int) int {
	return int(math.Floor(float64(defaultFactor * weight) / hr.averageWeight))
}

// create a hash ring
func (hr *HashRing) init(nodes []string, weights map[string]int) ([]*Node, error) {
	nodesCount := len(nodes)
	// check weights and initialize it
	if len(weights) == 0 {
		weights = make(map[string]int)
		for _, node := range nodes {
			weights[node] = defaultWeight
		}
	}
	// calculate total weight
	totalWeight := 0
	for _, weight := range weights {
		totalWeight += weight
	}
	hr.averageWeight = float64(totalWeight) / float64(nodesCount)
	// generate ring
	var rNodes []*Node
	for _, node := range nodes {
		weight := weights[node]
		rNode := hr.addNode(node, weight)
		// initialization failed
		if rNode == nil {
			continue
		}
		rNodes = append(rNodes, rNode)
	}
	return rNodes, nil
}

// with a powerlaw-alike distribution where higher
// levels are less likely to be returned
func (hr *HashRing) randomLevel() int {
	const branching = 4
	level := 1
	for hr.rnd.Int() % branching == 0 {
		level++
	}
	if level > maxLevel {
		return maxLevel
	}
	return level
}

// add server to hash ring
func (hr *HashRing) insertNode(rNode *Node, hash uint32) error {
	// record the server that are inserted into the location per level
	node := hr.header
	update := make(map[int]*VNode)
	for i := hr.level - 1; i >= 0; i-- {
		for {
			if node.Forward[i] != nil && node.Forward[i].Hash < hash {
				node = node.Forward[i]
			} else {
				break
			}
		}
		update[i] = node
	}
	level := hr.randomLevel()
	// filling
	if level > hr.level {
		for i := hr.level; i < level; i++ {
			update[i] = hr.header
		}
		hr.level = level
	}
	// update skip list
	newNode := NewVNode(rNode, hash, level)
	for i := 0; i < level; i++ {
		newNode.Forward[i] = update[i].Forward[i]
		update[i].Forward[i] = newNode
	}
	hr.length++
	return nil
}

func (hr *HashRing) genKey(v ...string) []byte {
	b := bytes.Buffer{}
	var i int
	var length = len(v)
	for i = 0; i < length - 1; i++ {
		b.WriteString(v[i])
		b.WriteString("#")
	}
	b.WriteString(v[i])
	return b.Bytes()
}

func (hr *HashRing) addNode(nodeIpaddr string, weight int) *Node {
	// check whether exist or not
	nodeIp := strings.Split(nodeIpaddr, ":")[0]
	if hr.filter.Lookup([]byte(nodeIp)) {
		return nil
	}
	rNode := NewNode(nodeIpaddr, weight)
	if rNode == nil {
		return nil
	}
	// add to bloom filter
	hr.filter.Insert([]byte(nodeIp))
	vNodeCount := hr.getVNodeCount(weight)
	// four virtual nodes per group
	for i := 0; i < vNodeCount / 4; i++ {
		hashKey := Md5Hash(hr.genKey(nodeIpaddr, string(i)))
		for j := 0; j < 4; j++ {
			key := KemataHash(hashKey, j)
			hr.insertNode(rNode, key)
		}
	}
	return rNode
}

func (hr *HashRing) removeNode(hash uint32) {
	node := hr.header
	update := make(map[int]*VNode)
	for i := hr.level - 1; i >= 0; i-- {
		for node.Forward[i] != nil && node.Forward[i].Hash < hash {
			node = node.Forward[i]
		}
		update[i] = node
	}
	// remove current server
	for i := hr.level - 1; i >= 0; i-- {
		update[i].Forward[i] = update[i].Forward[i].Forward[i]
	}
	// remove invalid level
	for hr.level > 1 && hr.header.Forward[hr.level - 1] == nil {
		hr.level--
	}
	hr.length--
}

// delete virtual node
func (hr *HashRing) deleteVnode(node *VNode) {
	if node == nil {
		return
	}
	hr.removeNode(node.Hash)
}

func (hr *HashRing) deleteNode(nodeIpaddr string, weight int) {
	// if not exist, return at once
	nodeIp := strings.Split(nodeIpaddr, ":")[0]
	if !hr.filter.Lookup([]byte(nodeIp)) {
		return
	}
	// delete node from cuckoo filter
	hr.filter.Delete([]byte(nodeIp))

	// do delete
	vNodeCount := hr.getVNodeCount(weight)
	for i := 0; i < vNodeCount / 4; i++ {
		hashKey := Md5Hash(hr.genKey(nodeIpaddr, string(i)))
		for j := 0; j < 4; j++ {
			key := KemataHash(hashKey, j)
			hr.removeNode(key)
		}
	}
}

// find a proper server for data
func (hr *HashRing) findProperNode(hashKey uint32) (*VNode, error) {
	// the hashKey is hash value of data, instead of node's
	if hr.header.Forward[0] == nil {
		return nil, errors.New("no node in this hash ring")
	}
	node := hr.header
	for i := hr.level - 1; i >= 0; i-- {
		for node.Forward[i] != nil && node.Forward[i].Hash < hashKey {
			node = node.Forward[i]
		}
	}
	// arrive the end, turn to head
	if node.Forward[0] == nil {
		node = hr.header.Forward[0]
	}
	return node, nil
}

// for debug
// iterate this ring
func (hr *HashRing) Iter(f func(*VNode)) {
	root := hr.header.Forward[0]
	for p := root; p != nil; p = p.Forward[0] {
		f(p)
	}
}

// return the first virtual node in this ring
func (hr *HashRing) First() *VNode { return hr.header.Forward[0] }