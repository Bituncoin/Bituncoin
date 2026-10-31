package api

import (
	"fmt"
	"net"
	"testing"
)

// freePort returns a free TCP port on localhost. It briefly binds the port
// and then closes the listener so the port is available for subsequent use.
func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to find free port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()
	return port
}

// newTestNode creates a Node with a guaranteed-free port pair (p, p+1)
// so that neither the HTTP listener nor the P2P listener requires root privileges.
func newTestNode(t *testing.T) (*Node, int) {
	t.Helper()
	// Find two consecutive free ports: p (API) and p+1 (P2P).
	for attempts := 0; attempts < 10; attempts++ {
		p := freePort(t)
		// Verify that p+1 is also available before committing.
		l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", p+1))
		if err != nil {
			continue
		}
		l.Close()

		node, err := NewNode("127.0.0.1", p)
		if err != nil {
			t.Fatalf("NewNode: %v", err)
		}
		return node, p
	}
	t.Fatal("could not find two consecutive free ports")
	return nil, 0
}

func TestNewNode(t *testing.T) {
	node, _ := newTestNode(t)

	if node == nil {
		t.Fatal("Expected non-nil node")
	}

	if node.p2pNetwork == nil {
		t.Error("Expected p2pNetwork to be initialized")
	}

	if node.IsRunning {
		t.Error("Expected node to not be running after creation")
	}
}

func TestNodeStartStop(t *testing.T) {
	node, _ := newTestNode(t)

	if err := node.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}

	if !node.IsRunning {
		t.Error("Expected node to be running after Start()")
	}

	if node.p2pNetwork != nil && !node.p2pNetwork.IsRunning() {
		t.Error("Expected p2pNetwork to be running after Start()")
	}

	if err := node.Stop(); err != nil {
		t.Fatalf("Stop: %v", err)
	}

	if node.IsRunning {
		t.Error("Expected node to not be running after Stop()")
	}

	if node.p2pNetwork != nil && node.p2pNetwork.IsRunning() {
		t.Error("Expected p2pNetwork to not be running after Stop()")
	}
}

func TestNodeStartAlreadyRunning(t *testing.T) {
	node, _ := newTestNode(t)

	if err := node.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	defer node.Stop()

	if err := node.Start(); err == nil {
		t.Error("Expected error when starting an already running node")
	}
}

func TestNodeStopNotRunning(t *testing.T) {
	node, _ := newTestNode(t)

	if err := node.Stop(); err == nil {
		t.Error("Expected error when stopping a node that is not running")
	}
}
