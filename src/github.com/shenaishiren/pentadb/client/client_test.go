// This is test file for client.go

package client

import (
	"testing"
)

//func TestNewClient_NoEnoughNodes(t *testing.T) {
//	var nodes = []string{"127.0.0.1:5000"}
//
//	if _, err := NewClient(nodes, nil, 3); err != nil {
//		t.Errorf(err.Error())
//	}
//}

func TestNewClient(t *testing.T) {
	var nodes = []string{
		"10.20.204.75:4567",
		"10.19.126.55:4567",
	}
	client, err := NewClient(nodes, nil, 1)
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer client.Close()
	if len(client.nodes) != len(nodes) {
		t.Error("wrong node number")
	}
	client.Put([]byte("p"), []byte("v"))
	if value := client.Get([]byte("p")); value == nil {
		t.Error("wrong get")
	} else {
		LOG.Debug("value: ", value)
	}
}




