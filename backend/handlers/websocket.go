package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"forum/models"
	"forum/utils"
	"forum/utils/middlewares"

	"github.com/gorilla/websocket"
)

var (
	// upgrader upgrades HTTP connections to WebSocket connections
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	// clients stores active WebSocket connections mapped by user ID
	clients   = make(map[int]*websocket.Conn)
	// clientsMu protects concurrent access to clients map
	clientsMu sync.RWMutex
)

// WebSocketHandler handles WebSocket connection upgrades and message routing
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket Upgrade Error:", err)
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}
	fmt.Println("WebSocket Connected")
	defer conn.Close()

	// Get authenticated user from context
	user, ok := r.Context().Value(middlewares.UserIDKey).(models.User)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Register client connection
	userID := user.ID
	clientsMu.Lock()
	clients[userID] = conn
	clientsMu.Unlock()

	// Broadcast user online status to all clients
	broadcastUserStatus()
	defer removeClient(userID)

	// Listen for incoming messages
	for {
		var msg models.WSMessage
		if err := conn.ReadJSON(&msg); err != nil {
			fmt.Println("WebSocket Read Error:", err)
			break
		}

		msg.SenderID = userID
		handleMessage(msg, conn)
	}
}

// ErrorMessage sends an error message to the client
func ErrorMessage(senderID int, conn *websocket.Conn) {
	conn.WriteJSON(models.WSMessage{
		Type: "error",
	})
}

// handleMessage processes different types of WebSocket messages
func handleMessage(msg models.WSMessage, conn *websocket.Conn) {
	var err error
	switch msg.Type {
	case "new_message":
		// Validate message content
		msg.Message.Content = strings.TrimSpace(msg.Message.Content)
		if len(msg.Message.Content) == 0 || len(msg.Message.Content) > 500 {
			ErrorMessage(msg.SenderID, conn)
			return
		}

		// Store message in database
		if err := msg.Message.StoreMessage(); err != nil {
			fmt.Println("Message Store Error:", err)
			ErrorMessage(msg.SenderID, conn)
			return
		}

		// Send message to sender and receiver
		distributeMessage(msg)

	case "read":
		// Mark messages as read
		msg.Message.UpdateRead()

	case "users":
		// Send list of users with online status
		msg.Type = "users"
		clientsMu.RLock()
		msg.Data = getClientIDs()
		clientsMu.RUnlock()
		msg.Members, err = GetUsers(msg.SenderID)
		if err != nil {
			fmt.Println("GetUsers Error:", err)
			return
		}
		conn.WriteJSON(msg)
		
	case "typing":
		// Broadcast typing indicator
		distributeMessage(msg)
	}
}

// distributeMessage sends a message to both sender and receiver
func distributeMessage(msg models.WSMessage) {
	var err error
	clientsMu.RLock()
	defer clientsMu.RUnlock()

	// Send to both sender and receiver
	for _, recipientID := range []int{msg.SenderID, msg.ReceiverID} {
		if conn, ok := clients[recipientID]; ok {
			msg.Members, err = GetUsers(recipientID)
			if err != nil {
				fmt.Println("GetUsers Error:", err)
				return
			}
			msg.Data = getClientIDs()
			if err := conn.WriteJSON(msg); err != nil {
				fmt.Println("WebSocket Write Error:", err)
			}
		}
	}
}

// broadcastUserStatus sends online user list to all connected clients
func broadcastUserStatus() {
	var err error
	clientsMu.RLock()
	defer clientsMu.RUnlock()

	msg := models.WSMessage{
		Type: "status_update",
		Data: getClientIDs(),
	}

	// Send to all connected clients
	for id, client := range clients {
		msg.Members, err = GetUsers(id)
		if err != nil {
			fmt.Println("GetUsers Error:", err)
			return
		}
		if err := client.WriteJSON(msg); err != nil {
			fmt.Println("WebSocket Broadcast Error:", err)
		}
	}
}

// removeClient closes and removes a client connection
func removeClient(userID int) {
	clientsMu.Lock()
	if conn, ok := clients[userID]; ok {
		conn.Close()
		delete(clients, userID)
	}
	clientsMu.Unlock()

	// Notify other clients that user went offline
	broadcastUserStatus()
}

// getClientIDs returns a list of all connected user IDs
func getClientIDs() []int {
	keys := make([]int, 0, len(clients))
	for key := range clients {
		keys = append(keys, key)
	}
	return keys
}
