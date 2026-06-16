package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"secure-p2p-backend/internal/crypto"
	"secure-p2p-backend/internal/models"
	"secure-p2p-backend/internal/signaling"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for now - configure appropriately for production
		return true
	},
}

// Server represents the main server with all handlers
type Server struct {
	userStore   map[string]*models.User
	sessionStore map[string]*models.Session
	roomManager *signaling.RoomManager
	mu          sync.RWMutex
	connections map[string]*websocket.Conn // connectionID -> websocket.Conn
}

// NewServer creates a new server instance
func NewServer() *Server {
	return &Server{
		userStore:    make(map[string]*models.User),
		sessionStore: make(map[string]*models.Session),
		roomManager:  signaling.NewRoomManager(10), // Max 10 peers per room
		connections:  make(map[string]*websocket.Conn),
	}
}

// Register handles user registration
func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		s.sendError(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if user already exists
	if _, exists := s.userStore[req.Username]; exists {
		s.sendError(w, "Username already exists", http.StatusConflict)
		return
	}

	// Hash password securely
	passwordHash, err := crypto.HashPassword(req.Password)
	if err != nil {
		s.sendError(w, "Internal error", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now().Unix(),
	}

	s.userStore[req.Username] = user

	s.sendSuccess(w, map[string]string{
		"user_id":  user.ID,
		"username": user.Username,
	})
}

// Login handles user authentication
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.sendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	s.mu.RLock()
	user, exists := s.userStore[req.Username]
	s.mu.RUnlock()

	if !exists {
		s.sendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Verify password securely
	if !crypto.VerifyPassword(req.Password, user.PasswordHash) {
		s.sendError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate secure session token
	token, err := crypto.GenerateSecureToken(32)
	if err != nil {
		s.sendError(w, "Internal error", http.StatusInternalServerError)
		return
	}

	session := &models.Session{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	s.mu.Lock()
	s.sessionStore[token] = session
	s.mu.Unlock()

	s.sendSuccess(w, map[string]interface{}{
		"token":      token,
		"user_id":    user.ID,
		"username":   user.Username,
		"expires_at": session.ExpiresAt,
	})
}

// WebSocket handles WebSocket connections for P2P signaling
func (s *Server) WebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	connectionID := uuid.New().String()
	
	s.mu.Lock()
	s.connections[connectionID] = conn
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.connections, connectionID)
		s.mu.Unlock()
		conn.Close()
	}()

	// Handle messages
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Parse and route message
		go s.handleWSMessage(conn, connectionID, message)
	}
}

func (s *Server) handleWSMessage(conn *websocket.Conn, connectionID string, message []byte) {
	var wsMsg models.WSMessage
	if err := json.Unmarshal(message, &wsMsg); err != nil {
		s.sendWSError(conn, "Invalid message format")
		return
	}

	// Verify message integrity (HMAC signature check would go here in production)
	// For now, we'll process the message type

	switch wsMsg.Type {
	case "join_room":
		s.handleJoinRoom(conn, connectionID, wsMsg.Payload)
	case "leave_room":
		s.handleLeaveRoom(conn, connectionID, wsMsg.Payload)
	case "offer":
		s.handleOffer(conn, connectionID, wsMsg.Payload)
	case "answer":
		s.handleAnswer(conn, connectionID, wsMsg.Payload)
	case "candidate":
		s.handleCandidate(conn, connectionID, wsMsg.Payload)
	default:
		s.sendWSError(conn, "Unknown message type")
	}
}

func (s *Server) handleJoinRoom(conn *websocket.Conn, connectionID string, payload string) {
	var msg struct {
		RoomID string `json:"room_id"`
		UserID string `json:"user_id"`
	}

	if err := json.Unmarshal([]byte(payload), &msg); err != nil {
		s.sendWSError(conn, "Invalid join_room payload")
		return
	}

	// Add peer to room
	peer, err := s.roomManager.AddPeer(msg.RoomID, uuid.New().String(), msg.UserID, connectionID)
	if err != nil {
		s.sendWSError(conn, err.Error())
		return
	}

	// Get other peers to notify them
	otherPeers := s.roomManager.GetOtherPeers(msg.RoomID, peer.ID)

	response := map[string]interface{}{
		"type":    "room_joined",
		"room_id": msg.RoomID,
		"peer_id": peer.ID,
		"peers":   otherPeers,
	}

	s.sendWSMessage(conn, response)

	// Notify other peers about new participant
	for _, otherPeer := range otherPeers {
		if otherConn, exists := s.connections[otherPeer.ConnectionID]; exists {
			notification := map[string]interface{}{
				"type":      "peer_joined",
				"room_id":   msg.RoomID,
				"new_peer":  peer,
			}
			s.sendWSMessage(otherConn, notification)
		}
	}
}

func (s *Server) handleLeaveRoom(conn *websocket.Conn, connectionID string, payload string) {
	var msg struct {
		RoomID string `json:"room_id"`
		PeerID string `json:"peer_id"`
	}

	if err := json.Unmarshal([]byte(payload), &msg); err != nil {
		s.sendWSError(conn, "Invalid leave_room payload")
		return
	}

	// Check if room exists before removing peer
	if _, exists := s.roomManager.GetRoom(msg.RoomID); !exists {
		s.sendWSError(conn, "Room not found")
		return
	}

	// Remove peer
	s.roomManager.RemovePeer(msg.RoomID, msg.PeerID)

	// Notify other peers
	peers := s.roomManager.GetOtherPeers(msg.RoomID, msg.PeerID)
	for _, peer := range peers {
		if peerConn, exists := s.connections[peer.ConnectionID]; exists {
			notification := map[string]interface{}{
				"type":      "peer_left",
				"room_id":   msg.RoomID,
				"peer_id":   msg.PeerID,
			}
			s.sendWSMessage(peerConn, notification)
		}
	}

	response := map[string]interface{}{
		"type":    "room_left",
		"room_id": msg.RoomID,
	}
	s.sendWSMessage(conn, response)
}

func (s *Server) handleOffer(conn *websocket.Conn, connectionID string, payload string) {
	var msg models.SignalingMessage
	if err := json.Unmarshal([]byte(payload), &msg); err != nil {
		s.sendWSError(conn, "Invalid offer payload")
		return
	}

	// Find target peer and forward offer
	if targetPeer, exists := s.roomManager.GetPeer(msg.RoomID, msg.TargetID); exists {
		if targetConn, exists := s.connections[targetPeer.ConnectionID]; exists {
			// Update SDP offer
			s.roomManager.UpdatePeerSDPOffer(msg.RoomID, msg.SenderID, msg.SDP)
			
			response := map[string]interface{}{
				"type":      "offer",
				"room_id":   msg.RoomID,
				"sender_id": msg.SenderID,
				"sdp":       msg.SDP,
			}
			s.sendWSMessage(targetConn, response)
		}
	}
}

func (s *Server) handleAnswer(conn *websocket.Conn, connectionID string, payload string) {
	var msg models.SignalingMessage
	if err := json.Unmarshal([]byte(payload), &msg); err != nil {
		s.sendWSError(conn, "Invalid answer payload")
		return
	}

	// Find target peer and forward answer
	if targetPeer, exists := s.roomManager.GetPeer(msg.RoomID, msg.TargetID); exists {
		if targetConn, exists := s.connections[targetPeer.ConnectionID]; exists {
			// Update SDP answer
			s.roomManager.UpdatePeerSDPAnswer(msg.RoomID, msg.SenderID, msg.SDP)
			
			response := map[string]interface{}{
				"type":      "answer",
				"room_id":   msg.RoomID,
				"sender_id": msg.SenderID,
				"sdp":       msg.SDP,
			}
			s.sendWSMessage(targetConn, response)
		}
	}
}

func (s *Server) handleCandidate(conn *websocket.Conn, connectionID string, payload string) {
	var msg models.SignalingMessage
	if err := json.Unmarshal([]byte(payload), &msg); err != nil {
		s.sendWSError(conn, "Invalid candidate payload")
		return
	}

	// Add ICE candidate
	s.roomManager.AddICECandidate(msg.RoomID, msg.SenderID, string(msg.Candidate.(string)))

	// Forward to target peer
	if targetPeer, exists := s.roomManager.GetPeer(msg.RoomID, msg.TargetID); exists {
		if targetConn, exists := s.connections[targetPeer.ConnectionID]; exists {
			response := map[string]interface{}{
				"type":      "candidate",
				"room_id":   msg.RoomID,
				"sender_id": msg.SenderID,
				"candidate": msg.Candidate,
			}
			s.sendWSMessage(targetConn, response)
		}
	}
}

func (s *Server) sendWSMessage(conn *websocket.Conn, data interface{}) {
	message, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling WS message: %v", err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Printf("Error sending WS message: %v", err)
	}
}

func (s *Server) sendWSError(conn *websocket.Conn, errorMsg string) {
	s.sendWSMessage(conn, map[string]interface{}{
		"type":  "error",
		"error": errorMsg,
	})
}

// CreateRoom handles room creation via REST API
func (s *Server) CreateRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Authenticate user (simplified - in production, verify session token)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		s.sendError(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	roomID := uuid.New().String()
	room := s.roomManager.CreateRoom(roomID, userID)

	s.sendSuccess(w, map[string]interface{}{
		"room_id": room.ID,
		"creator_id": room.CreatorID,
		"created_at": room.CreatedAt,
	})
}

// ListRooms returns all active rooms
func (s *Server) ListRooms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rooms := s.roomManager.ListRooms()
	s.sendSuccess(w, map[string]interface{}{
		"rooms": rooms,
	})
}

// Helper methods
func (s *Server) sendSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.APIResponse{
		Success: true,
		Data:    data,
	})
}

func (s *Server) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.APIResponse{
		Success: false,
		Error:   message,
	})
}

// SetupRoutes configures HTTP routes
func (s *Server) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/register", s.Register)
	mux.HandleFunc("/api/login", s.Login)
	mux.HandleFunc("/api/room", s.CreateRoom)
	mux.HandleFunc("/api/rooms", s.ListRooms)
	mux.HandleFunc("/ws", s.WebSocket)
}
