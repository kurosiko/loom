package signaling

import (
	"encoding/json"
	"sync"
	"time"

	"secure-p2p-backend/internal/models"
)

// RoomManager manages P2P video call rooms
type RoomManager struct {
	rooms   map[string]*models.Room
	mu      sync.RWMutex
	maxPeers int
}

// NewRoomManager creates a new room manager
func NewRoomManager(maxPeers int) *RoomManager {
	return &RoomManager{
		rooms:    make(map[string]*models.Room),
		maxPeers: maxPeers,
	}
}

// CreateRoom creates a new room with the given creator
func (rm *RoomManager) CreateRoom(roomID, creatorID string) *models.Room {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room := &models.Room{
		ID:        roomID,
		CreatorID: creatorID,
		Peers:     make(map[string]*models.Peer),
		CreatedAt: time.Now().Unix(),
		Metadata:  make(map[string]string),
	}

	rm.rooms[roomID] = room
	return room
}

// GetRoom retrieves a room by ID
func (rm *RoomManager) GetRoom(roomID string) (*models.Room, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	room, exists := rm.rooms[roomID]
	return room, exists
}

// AddPeer adds a peer to a room
func (rm *RoomManager) AddPeer(roomID, peerID, userID, connectionID string) (*models.Peer, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room, exists := rm.rooms[roomID]
	if !exists {
		return nil, ErrRoomNotFound
	}

	if len(room.Peers) >= rm.maxPeers {
		return nil, ErrRoomFull
	}

	peer := &models.Peer{
		ID:           peerID,
		UserID:       userID,
		ConnectionID: connectionID,
		Connected:    true,
		JoinedAt:     time.Now().Unix(),
	}

	room.Peers[peerID] = peer
	return peer, nil
}

// RemovePeer removes a peer from a room
func (rm *RoomManager) RemovePeer(roomID, peerID string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room, exists := rm.rooms[roomID]
	if !exists {
		return ErrRoomNotFound
	}

	delete(room.Peers, peerID)

	// Clean up empty rooms
	if len(room.Peers) == 0 {
		delete(rm.rooms, roomID)
	}

	return nil
}

// GetPeer retrieves a peer from a room
func (rm *RoomManager) GetPeer(roomID, peerID string) (*models.Peer, bool) {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	room, exists := rm.rooms[roomID]
	if !exists {
		return nil, false
	}

	peer, exists := room.Peers[peerID]
	return peer, exists
}

// UpdatePeerSDPOffer updates a peer's SDP offer
func (rm *RoomManager) UpdatePeerSDPOffer(roomID, peerID, sdpOffer string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room, exists := rm.rooms[roomID]
	if !exists {
		return ErrRoomNotFound
	}

	peer, exists := room.Peers[peerID]
	if !exists {
		return ErrPeerNotFound
	}

	peer.SDPOffer = sdpOffer
	return nil
}

// UpdatePeerSDPAnswer updates a peer's SDP answer
func (rm *RoomManager) UpdatePeerSDPAnswer(roomID, peerID, sdpAnswer string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room, exists := rm.rooms[roomID]
	if !exists {
		return ErrRoomNotFound
	}

	peer, exists := room.Peers[peerID]
	if !exists {
		return ErrPeerNotFound
	}

	peer.SDPAnswer = sdpAnswer
	return nil
}

// AddICECandidate adds an ICE candidate to a peer
func (rm *RoomManager) AddICECandidate(roomID, peerID, candidate string) error {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	room, exists := rm.rooms[roomID]
	if !exists {
		return ErrRoomNotFound
	}

	peer, exists := room.Peers[peerID]
	if !exists {
		return ErrPeerNotFound
	}

	peer.ICECandidates = append(peer.ICECandidates, candidate)
	return nil
}

// GetOtherPeers returns all peers in a room except the specified one
func (rm *RoomManager) GetOtherPeers(roomID, excludePeerID string) []*models.Peer {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	room, exists := rm.rooms[roomID]
	if !exists {
		return nil
	}

	peers := make([]*models.Peer, 0)
	for id, peer := range room.Peers {
		if id != excludePeerID {
			peers = append(peers, peer)
		}
	}

	return peers
}

// DeleteRoom deletes a room
func (rm *RoomManager) DeleteRoom(roomID string) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	delete(rm.rooms, roomID)
}

// ListRooms returns all active rooms
func (rm *RoomManager) ListRooms() []*models.Room {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	rooms := make([]*models.Room, 0, len(rm.rooms))
	for _, room := range rm.rooms {
		rooms = append(rooms, room)
	}

	return rooms
}

// MarshalRoom safely marshals a room to JSON (excluding sensitive data)
func (rm *RoomManager) MarshalRoom(room *models.Room) ([]byte, error) {
	// Create a safe copy without internal connection details
	safeRoom := struct {
		ID        string            `json:"id"`
		CreatorID string            `json:"creator_id"`
		PeerCount int               `json:"peer_count"`
		CreatedAt int64             `json:"created_at"`
		Metadata  map[string]string `json:"metadata,omitempty"`
	}{
		ID:        room.ID,
		CreatorID: room.CreatorID,
		PeerCount: len(room.Peers),
		CreatedAt: room.CreatedAt,
		Metadata:  room.Metadata,
	}

	return json.Marshal(safeRoom)
}

// Errors
var (
	ErrRoomNotFound  = &SignalingError{"room not found"}
	ErrRoomFull      = &SignalingError{"room is full"}
	ErrPeerNotFound  = &SignalingError{"peer not found"}
)

// SignalingError represents a signaling error
type SignalingError struct {
	Message string
}

func (e *SignalingError) Error() string {
	return e.Message
}
