// This is test file for client.go

package client

import (
	"testing"
)

//func TestNewClient_NoEnoughNodes(t *testing.T) {
//	var nodes = []string{"127.0.0.1:5000"}
//	var weights = map[string]int{"127.0.0.1:5000": 1}
//	if _, err := NewClient(nodes, weights, 3); err != nil {
//		t.Errorf(err.Error())
//	}
//}

func TestNewClient_WrongReplicas(t *testing.T) {
	var nodes = []string{"127.0.0.1:5000"}
	var weights = map[string]int{"127.0.0.1:5000": 1}
	if _, err := NewClient(nodes, weights, 3); err != nil {
		t.Errorf(err.Error())
	}
}

//func TestNewClient_NodeUnreachable(t *testing.T) {
//	var nodes = []string{"127.0.0.1:5001"}
//	var weights = map[string]int{"127.0.0.1:5001": 1}
//	if _, err := NewClient(nodes, weights, 1); err != nil {
//		t.Errorf(err.Error())
//	}
//}




