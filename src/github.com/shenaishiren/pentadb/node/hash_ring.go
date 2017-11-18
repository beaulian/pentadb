// Contains the interface and implementation of Hash Ring
// implemented by skip list

package node

import (
	"errors"
	"math"
	"math/rand"

	"github.com/shenaishiren/pentadb/utils"
)

const (
	defaultFactor = 40
	defaultWeight = 1
	maxLevel = 32
)


type HashRing struct {
	rnd *rand.Rand

	header *Node    // point to node skiplist
	level int
	maxHeight int
	length int      // number of nodes
	totalWeight int // total weight of nodes in hash ring
}

func NewHashRing() *HashRing {
	return &HashRing{
		rnd:          rand.New(rand.NewSource(0xdeadbeef)),
		level:        1,
		maxHeight:    1,
		length:       0,
		totalWeight:  0,
		header:       NewNode("nil",0, maxLevel),
	}
}

// Get the count of virtual nodes
func (hr *HashRing) getVNodeCount(weight int, averageWeight float64) int {
	return int(math.Floor(float64(defaultFactor * weight) / averageWeight))
}

// create a hash ring
func (hr *HashRing) init(nodes []string, weights map[string]int) error {
	// check nodes' count
	nodesCount := len(nodes)
	if nodesCount < 3 {
		return errors.New("nodes number must be greater than 3")
	}
	// check weights and initialize it
	if weights == nil {
		weights := make(map[string]int)
		for _, node := range nodes {
			weights[node] = defaultWeight
		}
	}
	// calculate total weight
	totalWeight := 0
	for _, weight := range weights {
		totalWeight += weight
	}
	averageWeight := float64(totalWeight / nodesCount)
	// generate ring
	for _, node := range nodes {
		weight := weights[node]
		// get virtual nodes' count
		vNodeCount := hr.getVNodeCount(weight, averageWeight)
		hr.addNode(node, vNodeCount)
	}
	return nil
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


// add node to hash ring
func (hr *HashRing) insertNode(nodeName string, hash uint32) error {
	// record the node that are inserted into the location per level
	node := hr.header
	update := make(map[int]*Node)
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
	newNode := NewNode(nodeName, hash, level)
	for i := 0; i < level; i++ {
		newNode.Forward[i] = update[i].Forward[i]
		update[i].Forward[i] = newNode
	}
	hr.length++

	return nil
}

func (hr *HashRing) addNode(nodeName string, vNodeCount int) error {
	// four virtual nodes per group
	for i := 0; i < vNodeCount / 4; i++ {
		for j := 0; j < 4; j++ {
			key := utils.KemataHash(nodeName, j)
			hr.insertNode(nodeName, key)
		}
	}
	return nil
}

func (hr *HashRing) removeNode(hash uint32) error {
	node := hr.header
	update := make(map[int]*Node)
	for i := hr.level - 1; i >= 0; i-- {
		for node.Forward[i].Hash < hash {
			node = node.Forward[i]
		}
		update[i] = node
	}
	// remove current node
	for i := hr.level - 1; i >= 0; i-- {
		update[i].Forward[i] = update[i].Forward[i].Forward[i]
	}
	// remove invalid level
	for hr.level > 1 && hr.header.Forward[hr.level - 1] == nil {
		hr.level--
	}
	hr.length--

	return nil
}

func (hr *HashRing) deleteNode(node *Node) error {
	hr.removeNode(node.Hash)

	return  nil
}

func (hr *HashRing) deleteNodeByName(nodeName string, vNodeCount int) error {
	for i := 0; i < vNodeCount / 4; i++ {
		for j := 0; j < 4; j++ {
			key := utils.KemataHash(nodeName, j)
			hr.removeNode(key)
		}
	}

	return nil
}

// find a proper node for data
func (hr *HashRing) findProperNode(hashKey uint32) (*Node, error) {
	// the hashKey is hash value of data, instead of node
	node := hr.header
	for i := hr.level - 1; i >= 0; i-- {
		for node.Forward[i].Hash < hashKey {
			node = node.Forward[i]
		}
	}
	// arrive the end
	if node.Forward[0] == nil {
		node = hr.header.Forward[0]
	}
	return node, nil
}