package models

// User represents a user in the system
type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"` // Never expose password hash in JSON
	PublicKey    string `json:"public_key,omitempty"`
	CreatedAt    int64  `json:"created_at"`
}

// Session represents an authenticated session
type Session struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

// Room represents a P2P video call room
type Room struct {
	ID        string            `json:"id"`
	CreatorID string            `json:"creator_id"`
	Peers     map[string]*Peer  `json:"peers"`
	CreatedAt int64             `json:"created_at"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// Peer represents a participant in a room
type Peer struct {
	ID           string   `json:"id"`
	UserID       string   `json:"user_id"`
	ConnectionID string   `json:"connection_id"`
	SDPOffer     string   `json:"sdp_offer,omitempty"`
	SDPAnswer    string   `json:"sdp_answer,omitempty"`
	ICECandidates []string `json:"ice_candidates,omitempty"`
	Connected    bool     `json:"connected"`
	JoinedAt     int64    `json:"joined_at"`
}

// SignalingMessage represents a WebRTC signaling message
type SignalingMessage struct {
	Type      string      `json:"type"` // "offer", "answer", "candidate"
	RoomID    string      `json:"room_id"`
	SenderID  string      `json:"sender_id"`
	TargetID  string      `json:"target_id,omitempty"`
	SDP       string      `json:"sdp,omitempty"`
	Candidate interface{} `json:"candidate,omitempty"`
}

// WSMessage represents a WebSocket message with encryption support
type WSMessage struct {
	Type        string `json:"type"`
	Payload     string `json:"payload,omitempty"` // Encrypted payload
	Nonce       string `json:"nonce,omitempty"`   // Nonce for decryption
	Signature   string `json:"signature,omitempty"` // HMAC signature for integrity
	Timestamp   int64  `json:"timestamp"`
}

// EncryptedData represents encrypted data with metadata
type EncryptedData struct {
	Ciphertext string `json:"ciphertext"`
	Nonce      string `json:"nonce"`
	Tag        string `json:"tag"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CreateRoomRequest represents a room creation request
type CreateRoomRequest struct {
	Metadata map[string]string `json:"metadata,omitempty"`
}
