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
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true }, // Allow all origins
	}
	clients   = make(map[int]*websocket.Conn)
	clientsMu sync.RWMutex
)

// WebSocketHandler upgrades HTTP connections to WebSocket and manages real-time messaging.
// It handles client registration, message routing, and broadcasts user status updates.
// The connection is maintained until the client disconnects or an error occurs.
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket Upgrade Error:", err)
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}
	fmt.Println("WebSocket Connected")
	defer conn.Close()

	user, ok := r.Context().Value(middlewares.UserIDKey).(models.User)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID := user.ID
	clientsMu.Lock()
	clients[userID] = conn
	clientsMu.Unlock()

	broadcastUserStatus()
	defer removeClient(userID)

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

// ErrorMessage sends an error message to a WebSocket client.
// It's used to notify clients when their message couldn't be processed.
func ErrorMessage(senderID int, conn *websocket.Conn) {
	conn.WriteJSON(models.WSMessage{
		Type: "error",
	})
}

// handleMessage processes incoming WebSocket messages based on their type.
// Supported types: "new_message" (chat), "read" (mark as read), "users" (get user list), "typing" (typing indicator).
// It validates message content and routes them to appropriate handlers.
func handleMessage(msg models.WSMessage, conn *websocket.Conn) {
	var err error
	switch msg.Type {
	case "new_message":
		msg.Message.Content = strings.TrimSpace(msg.Message.Content)
		if len(msg.Message.Content) == 0 || len(msg.Message.Content) > 500 {
			ErrorMessage(msg.SenderID, conn)
			return
		}

		if err := msg.Message.StoreMessage(); err != nil {
			fmt.Println("Message Store Error:", err)
			ErrorMessage(msg.SenderID, conn)
			return
		}

		distributeMessage(msg)

	case "read":
		msg.Message.UpdateRead()

	case "users":
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
		distributeMessage(msg)
	}
}

// distributeMessage sends a message to both sender and receiver if they're connected.
// It updates the user list and member information for each recipient.
func distributeMessage(msg models.WSMessage) {
	var err error
	clientsMu.RLock()
	defer clientsMu.RUnlock()

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

// broadcastUserStatus sends the current list of online users to all connected clients.
// It's called when users connect or disconnect to keep all clients synchronized.
func broadcastUserStatus() {
	var err error
	clientsMu.RLock()
	defer clientsMu.RUnlock()

	msg := models.WSMessage{
		Type: "status_update",
		Data: getClientIDs(),
	}

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

// removeClient removes a user from the active clients map and closes their connection.
// It also broadcasts the updated user status to all remaining clients.
func removeClient(userID int) {
	clientsMu.Lock()
	if conn, ok := clients[userID]; ok {
		conn.Close()
		delete(clients, userID)
	}
	clientsMu.Unlock()

	broadcastUserStatus()
}

// getClientIDs returns a slice of all currently connected user IDs.
// This is used to inform clients about who is online.
func getClientIDs() []int {
	keys := make([]int, 0, len(clients))
	for key := range clients {
		keys = append(keys, key)
	}
	return keys
}
