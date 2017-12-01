package pentadb

import "testing"

func TestNewHashRing(t *testing.T) {
	hashRing := NewHashRing()
	nodes := []string{
		"127.0.0.1:5000",
		"127.0.0.1:5001",
		"127.0.0.1:5002",
	}
	weights := map[string]int{
		"127.0.0.1:5000": 1,
		"127.0.0.1:5001": 1,
		"127.0.0.1:5002": 1,
	}
	if rNodes, _ := hashRing.init(nodes, weights); len(rNodes) != 3 {
		t.Errorf("wrong nodes: %d, %s", len(rNodes), rNodes)
	}
	group := make(map[string]int)
	hashRing.Iter(func (v *VNode) {
		if _, ok := group[v.rNode.Ipaddr]; !ok {
			group[v.rNode.Ipaddr] = 0
		}
		group[v.rNode.Ipaddr]++
	})
	// test number
	for ip, num := range group {
		LOG.Debug(ip, ": ", num)
	}
	// test ring
	// need change code before test
	node, err := hashRing.findProperNode(KemataHash(Md5Hash([]byte("test")), 0))
	if err != nil {
		t.Error(err.Error())
	}
	if node != hashRing.First() {
		t.Error("hash ring is wrong!")
	}
	// test delete node
	node = hashRing.First()
	hashRing.deleteVnode(node)
	if hashRing.First() == node {
		t.Error("wrong delete function!")
	}
}
