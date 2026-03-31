package network

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewNetwork(t *testing.T) {
	n, err := NewNetwork("127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create network: %v", err)
	}

	if n.GetNodeID() == "" {
		t.Error("Expected non-empty node ID")
	}

	if n.IsRunning() {
		t.Error("Expected network to not be running before Start()")
	}
}

func TestNetworkStartStop(t *testing.T) {
	n, err := NewNetwork("127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create network: %v", err)
	}

	if err := n.Start(); err != nil {
		t.Fatalf("Failed to start network: %v", err)
	}

	if !n.IsRunning() {
		t.Error("Expected network to be running after Start()")
	}

	if err := n.Stop(); err != nil {
		t.Fatalf("Failed to stop network: %v", err)
	}

	if n.IsRunning() {
		t.Error("Expected network to not be running after Stop()")
	}
}

func TestNetworkStartAlreadyRunning(t *testing.T) {
	n, err := NewNetwork("127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create network: %v", err)
	}

	if err := n.Start(); err != nil {
		t.Fatalf("Failed to start network: %v", err)
	}
	defer n.Stop()

	if err := n.Start(); err == nil {
		t.Error("Expected error when starting an already running network")
	}
}

func TestNetworkStopNotRunning(t *testing.T) {
	n, err := NewNetwork("127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create network: %v", err)
	}

	if err := n.Stop(); err == nil {
		t.Error("Expected error when stopping a network that is not running")
	}
}

func TestNetworkGetStatus(t *testing.T) {
	n, err := NewNetwork("127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create network: %v", err)
	}

	status := n.GetStatus()

	if status["nodeId"] == "" {
		t.Error("Expected non-empty nodeId in status")
	}

	if status["version"] != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%v'", status["version"])
	}

	if status["isRunning"] != false {
		t.Error("Expected isRunning to be false before Start()")
	}

	if status["totalPeers"] != 0 {
		t.Errorf("Expected 0 total peers, got %v", status["totalPeers"])
	}

	if status["connectedPeers"] != 0 {
		t.Errorf("Expected 0 connected peers, got %v", status["connectedPeers"])
	}
}

func TestNetworkSetBlockHeight(t *testing.T) {
	n, err := NewNetwork("127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create network: %v", err)
	}

	n.SetBlockHeight(42)

	status := n.GetStatus()
	if status["blockHeight"] != 42 {
		t.Errorf("Expected blockHeight 42, got %v", status["blockHeight"])
	}
}

func TestNetworkGetPeersEmpty(t *testing.T) {
	n, err := NewNetwork("127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create network: %v", err)
	}

	peers := n.GetPeers()
	if len(peers) != 0 {
		t.Errorf("Expected 0 peers, got %d", len(peers))
	}
}

func TestNetworkGetPeerCountEmpty(t *testing.T) {
	n, err := NewNetwork("127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create network: %v", err)
	}

	if n.GetPeerCount() != 0 {
		t.Errorf("Expected peer count 0, got %d", n.GetPeerCount())
	}
}

func TestNetworkConnectDisconnect(t *testing.T) {
	// Create two network nodes
	server, err := NewNetwork("127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create server network: %v", err)
	}

	if err := server.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer server.Stop()

	serverAddr := server.listener.Addr().String()

	client, err := NewNetwork("127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create client network: %v", err)
	}

	peer, err := client.Connect(serverAddr)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}

	if peer.ID == "" {
		t.Error("Expected non-empty peer ID")
	}

	if peer.Address != serverAddr {
		t.Errorf("Expected peer address '%s', got '%s'", serverAddr, peer.Address)
	}

	if client.GetPeerCount() != 1 {
		t.Errorf("Expected 1 peer, got %d", client.GetPeerCount())
	}

	if err := client.Disconnect(peer.ID); err != nil {
		t.Fatalf("Failed to disconnect peer: %v", err)
	}

	if client.GetPeerCount() != 0 {
		t.Errorf("Expected 0 peers after disconnect, got %d", client.GetPeerCount())
	}
}

func TestNetworkDisconnectNonExistent(t *testing.T) {
	n, err := NewNetwork("127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create network: %v", err)
	}

	if err := n.Disconnect("nonexistent-peer-id"); err == nil {
		t.Error("Expected error when disconnecting non-existent peer")
	}
}

func TestNetworkConnectAlreadyConnected(t *testing.T) {
	server, err := NewNetwork("127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	if err := server.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer server.Stop()

	serverAddr := server.listener.Addr().String()

	client, err := NewNetwork("127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	peer, err := client.Connect(serverAddr)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Disconnect(peer.ID)

	// Wait briefly for the handshake to complete and status to become connected
	time.Sleep(50 * time.Millisecond)

	_, err = client.Connect(serverAddr)
	if err == nil {
		t.Error("Expected error when connecting to an already-connected peer")
	}
}

func TestNetworkBroadcast(t *testing.T) {
	server, err := NewNetwork("127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	if err := server.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer server.Stop()

	serverAddr := server.listener.Addr().String()

	client, err := NewNetwork("127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	peer, err := client.Connect(serverAddr)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Disconnect(peer.ID)

	// Wait for handshake to complete
	time.Sleep(50 * time.Millisecond)

	// Broadcast to all connected peers (should not error even if peer is still handshaking)
	err = client.Broadcast(MessageBlock, map[string]interface{}{"height": 1})
	// Broadcast may return nil or an error depending on connection state; we just ensure no panic
	_ = err
}

func TestNetworkSendToNonExistentPeer(t *testing.T) {
	n, err := NewNetwork("127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create network: %v", err)
	}

	err = n.Send("nonexistent", MessagePing, map[string]string{})
	if err == nil {
		t.Error("Expected error when sending to non-existent peer")
	}
}

func TestNetworkMessageHandler(t *testing.T) {
	received := make(chan *Message, 1)

	server, err := NewNetwork("127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	server.SetMessageHandler(func(peer *Peer, msg *Message) {
		if msg.Type == MessageTransaction {
			received <- msg
		}
	})

	if err := server.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer server.Stop()

	serverAddr := server.listener.Addr().String()

	client, err := NewNetwork("127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	peer, err := client.Connect(serverAddr)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Disconnect(peer.ID)

	// Wait for handshake
	time.Sleep(100 * time.Millisecond)

	// Send a transaction message from server side to trigger handler on the peer
	// Instead, broadcast from client to server
	_ = client.Broadcast(MessageTransaction, map[string]string{"txId": "tx123"})

	select {
	case msg := <-received:
		if msg.Type != MessageTransaction {
			t.Errorf("Expected MessageTransaction, got %s", msg.Type)
		}
		var payload map[string]string
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			t.Fatalf("Failed to unmarshal payload: %v", err)
		}
		if payload["txId"] != "tx123" {
			t.Errorf("Expected txId 'tx123', got '%s'", payload["txId"])
		}
	case <-time.After(2 * time.Second):
		t.Error("Timed out waiting for message")
	}
}

func TestNetworkGetConnectedPeers(t *testing.T) {
	server, err := NewNetwork("127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	if err := server.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer server.Stop()

	serverAddr := server.listener.Addr().String()

	client, err := NewNetwork("127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Before connecting, no connected peers
	if len(client.GetConnectedPeers()) != 0 {
		t.Errorf("Expected 0 connected peers, got %d", len(client.GetConnectedPeers()))
	}

	peer, err := client.Connect(serverAddr)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Disconnect(peer.ID)

	// Wait for handshake so status transitions to connected
	time.Sleep(100 * time.Millisecond)

	connected := client.GetConnectedPeers()
	if len(connected) != 1 {
		t.Errorf("Expected 1 connected peer, got %d", len(connected))
	}
}

func TestPeerStatusValues(t *testing.T) {
	if PeerStatusConnecting != "connecting" {
		t.Errorf("Unexpected PeerStatusConnecting value: %s", PeerStatusConnecting)
	}

	if PeerStatusConnected != "connected" {
		t.Errorf("Unexpected PeerStatusConnected value: %s", PeerStatusConnected)
	}

	if PeerStatusDisconnected != "disconnected" {
		t.Errorf("Unexpected PeerStatusDisconnected value: %s", PeerStatusDisconnected)
	}
}

func TestMessageTypeValues(t *testing.T) {
	types := []MessageType{
		MessageHandshake,
		MessageBlock,
		MessageTransaction,
		MessagePeerList,
		MessagePing,
		MessagePong,
	}

	for _, mt := range types {
		if mt == "" {
			t.Errorf("MessageType should not be empty")
		}
	}
}
