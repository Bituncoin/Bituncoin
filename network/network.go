package network

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

// MessageType represents the type of a network message
type MessageType string

const (
	MessageHandshake    MessageType = "handshake"
	MessageBlock        MessageType = "block"
	MessageTransaction  MessageType = "transaction"
	MessagePeerList     MessageType = "peer_list"
	MessagePing         MessageType = "ping"
	MessagePong         MessageType = "pong"
)

// PeerStatus represents the connection status of a peer
type PeerStatus string

const (
	PeerStatusConnecting   PeerStatus = "connecting"
	PeerStatusConnected    PeerStatus = "connected"
	PeerStatusDisconnected PeerStatus = "disconnected"
)

// Message represents a network protocol message
type Message struct {
	Type      MessageType     `json:"type"`
	From      string          `json:"from"`
	Timestamp int64           `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}

// HandshakePayload is the payload for a handshake message
type HandshakePayload struct {
	NodeID      string `json:"nodeId"`
	Version     string `json:"version"`
	Address     string `json:"address"`
	BlockHeight int    `json:"blockHeight"`
}

// Peer represents a connected peer node in the BTNG network
type Peer struct {
	ID          string     `json:"id"`
	Address     string     `json:"address"`
	Status      PeerStatus `json:"status"`
	Version     string     `json:"version"`
	BlockHeight int        `json:"blockHeight"`
	ConnectedAt time.Time  `json:"connectedAt"`
	LastSeen    time.Time  `json:"lastSeen"`
	conn        net.Conn
}

// Network manages the BTNG peer-to-peer network
type Network struct {
	nodeID      string
	listenAddr  string
	peers       map[string]*Peer
	listener    net.Listener
	isRunning   bool
	version     string
	blockHeight int
	onMessage   func(*Peer, *Message)
	mutex       sync.RWMutex
}

// NewNetwork creates a new BTNG network manager
func NewNetwork(listenAddr string) (*Network, error) {
	nodeID, err := generateNodeID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate node ID: %w", err)
	}

	return &Network{
		nodeID:     nodeID,
		listenAddr: listenAddr,
		peers:      make(map[string]*Peer),
		isRunning:  false,
		version:    "1.0.0",
	}, nil
}

// SetMessageHandler registers a callback invoked for every incoming message
func (n *Network) SetMessageHandler(handler func(*Peer, *Message)) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.onMessage = handler
}

// SetBlockHeight updates the locally advertised block height
func (n *Network) SetBlockHeight(height int) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.blockHeight = height
}

// Start begins listening for incoming peer connections
func (n *Network) Start() error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if n.isRunning {
		return errors.New("network already running")
	}

	ln, err := net.Listen("tcp", n.listenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", n.listenAddr, err)
	}

	n.listener = ln
	n.isRunning = true

	go n.acceptLoop(ln)

	return nil
}

// Stop shuts down the network listener and disconnects all peers
func (n *Network) Stop() error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if !n.isRunning {
		return errors.New("network not running")
	}

	n.isRunning = false

	if n.listener != nil {
		n.listener.Close()
		n.listener = nil
	}

	for _, peer := range n.peers {
		if peer.conn != nil {
			peer.conn.Close()
		}
		peer.Status = PeerStatusDisconnected
	}

	return nil
}

// Connect dials a remote peer and performs the BTNG handshake
func (n *Network) Connect(address string) (*Peer, error) {
	n.mutex.RLock()
	for _, p := range n.peers {
		if p.Address == address && p.Status == PeerStatusConnected {
			n.mutex.RUnlock()
			return nil, errors.New("already connected to peer")
		}
	}
	n.mutex.RUnlock()

	conn, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", address, err)
	}

	peer, err := n.registerPeer(conn, address)
	if err != nil {
		conn.Close()
		return nil, err
	}

	if err := n.sendHandshake(peer); err != nil {
		n.removePeer(peer.ID)
		return nil, fmt.Errorf("handshake failed: %w", err)
	}

	go n.readLoop(peer)

	return peer, nil
}

// Disconnect closes the connection to a peer by its ID
func (n *Network) Disconnect(peerID string) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	peer, exists := n.peers[peerID]
	if !exists {
		return errors.New("peer not found")
	}

	if peer.conn != nil {
		peer.conn.Close()
	}

	peer.Status = PeerStatusDisconnected
	delete(n.peers, peerID)

	return nil
}

// Broadcast sends a message to all connected peers
func (n *Network) Broadcast(msgType MessageType, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	msg := &Message{
		Type:      msgType,
		From:      n.nodeID,
		Timestamp: time.Now().Unix(),
		Payload:   json.RawMessage(data),
	}

	encoded, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	encoded = append(encoded, '\n')

	n.mutex.RLock()
	defer n.mutex.RUnlock()

	var errs []error
	for _, peer := range n.peers {
		if peer.Status == PeerStatusConnected && peer.conn != nil {
			if _, err := peer.conn.Write(encoded); err != nil {
				errs = append(errs, fmt.Errorf("peer %s: %w", peer.ID, err))
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("broadcast encountered %d error(s): %v", len(errs), errs)
	}

	return nil
}

// Send sends a message to a specific peer by ID
func (n *Network) Send(peerID string, msgType MessageType, payload interface{}) error {
	n.mutex.RLock()
	peer, exists := n.peers[peerID]
	n.mutex.RUnlock()

	if !exists {
		return errors.New("peer not found")
	}

	if peer.Status != PeerStatusConnected && peer.Status != PeerStatusConnecting {
		return errors.New("peer not connected")
	}

	if peer.conn == nil {
		return errors.New("peer connection is nil")
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	msg := &Message{
		Type:      msgType,
		From:      n.nodeID,
		Timestamp: time.Now().Unix(),
		Payload:   json.RawMessage(data),
	}

	encoded, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	encoded = append(encoded, '\n')

	_, err = peer.conn.Write(encoded)
	return err
}

// GetPeers returns a snapshot of all known peers
func (n *Network) GetPeers() []*Peer {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	peers := make([]*Peer, 0, len(n.peers))
	for _, p := range n.peers {
		peers = append(peers, p)
	}

	return peers
}

// GetConnectedPeers returns only peers with an active connection
func (n *Network) GetConnectedPeers() []*Peer {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	peers := make([]*Peer, 0)
	for _, p := range n.peers {
		if p.Status == PeerStatusConnected {
			peers = append(peers, p)
		}
	}

	return peers
}

// GetPeerCount returns the total number of known peers
func (n *Network) GetPeerCount() int {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return len(n.peers)
}

// GetNodeID returns this node's unique identifier
func (n *Network) GetNodeID() string {
	return n.nodeID
}

// IsRunning reports whether the network listener is active
func (n *Network) IsRunning() bool {
	n.mutex.RLock()
	defer n.mutex.RUnlock()
	return n.isRunning
}

// GetStatus returns a summary of the current network state
func (n *Network) GetStatus() map[string]interface{} {
	n.mutex.RLock()
	defer n.mutex.RUnlock()

	connected := 0
	for _, p := range n.peers {
		if p.Status == PeerStatusConnected {
			connected++
		}
	}

	return map[string]interface{}{
		"nodeId":         n.nodeID,
		"listenAddr":     n.listenAddr,
		"isRunning":      n.isRunning,
		"version":        n.version,
		"blockHeight":    n.blockHeight,
		"totalPeers":     len(n.peers),
		"connectedPeers": connected,
	}
}

// --- internal helpers ---

// acceptLoop waits for inbound TCP connections on the given listener
func (n *Network) acceptLoop(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			n.mutex.RLock()
			running := n.isRunning
			n.mutex.RUnlock()
			if !running {
				return
			}
			continue
		}

		peer, err := n.registerPeer(conn, conn.RemoteAddr().String())
		if err != nil {
			conn.Close()
			continue
		}

		// Greet the inbound peer so it can mark us as connected
		n.sendHandshake(peer) //nolint:errcheck

		go n.readLoop(peer)
	}
}

// readLoop reads newline-delimited JSON messages from a peer connection
func (n *Network) readLoop(peer *Peer) {
	decoder := json.NewDecoder(peer.conn)

	for {
		var msg Message
		if err := decoder.Decode(&msg); err != nil {
			n.mutex.Lock()
			if p, exists := n.peers[peer.ID]; exists {
				p.Status = PeerStatusDisconnected
			}
			n.mutex.Unlock()
			return
		}

		n.mutex.Lock()
		if p, exists := n.peers[peer.ID]; exists {
			p.LastSeen = time.Now()
		}
		n.mutex.Unlock()

		n.handleMessage(peer, &msg)
	}
}

// handleMessage dispatches an incoming message
func (n *Network) handleMessage(peer *Peer, msg *Message) {
	switch msg.Type {
	case MessageHandshake:
		var payload HandshakePayload
		if err := json.Unmarshal(msg.Payload, &payload); err == nil {
			n.mutex.Lock()
			if p, exists := n.peers[peer.ID]; exists {
				p.Version = payload.Version
				p.BlockHeight = payload.BlockHeight
				p.Status = PeerStatusConnected
			}
			n.mutex.Unlock()
		}

	case MessagePing:
		n.Send(peer.ID, MessagePong, map[string]int64{"timestamp": time.Now().Unix()}) //nolint:errcheck
	}

	// Forward to user-supplied handler
	n.mutex.RLock()
	handler := n.onMessage
	n.mutex.RUnlock()

	if handler != nil {
		handler(peer, msg)
	}
}

// sendHandshake sends a handshake message to a peer
func (n *Network) sendHandshake(peer *Peer) error {
	n.mutex.RLock()
	payload := HandshakePayload{
		NodeID:      n.nodeID,
		Version:     n.version,
		Address:     n.listenAddr,
		BlockHeight: n.blockHeight,
	}
	n.mutex.RUnlock()

	return n.Send(peer.ID, MessageHandshake, payload)
}

// registerPeer records a new peer connection
func (n *Network) registerPeer(conn net.Conn, address string) (*Peer, error) {
	peerID, err := generateNodeID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate peer ID: %w", err)
	}

	peer := &Peer{
		ID:          peerID,
		Address:     address,
		Status:      PeerStatusConnecting,
		ConnectedAt: time.Now(),
		LastSeen:    time.Now(),
		conn:        conn,
	}

	n.mutex.Lock()
	n.peers[peerID] = peer
	n.mutex.Unlock()

	return peer, nil
}

// removePeer removes a peer by ID and closes its connection
func (n *Network) removePeer(peerID string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if peer, exists := n.peers[peerID]; exists {
		if peer.conn != nil {
			peer.conn.Close()
		}
		delete(n.peers, peerID)
	}
}

// generateNodeID creates a cryptographically random 16-byte hex node identifier
func generateNodeID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
