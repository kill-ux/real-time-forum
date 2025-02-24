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

func handleMessage(msg models.WSMessage, conn *websocket.Conn) {
	var err error
	switch msg.Type {
	case "new_message":
		msg.Message.Content = strings.TrimSpace(msg.Message.Content)
		if len(msg.Message.Content) == 0 || len(msg.Message.Content) > 2000 {
			return
		}

		if err := msg.Message.StoreMessage(); err != nil {
			fmt.Println("Message Store Error:", err)
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
	}
}

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

func removeClient(userID int) {
	clientsMu.Lock()
	if conn, ok := clients[userID]; ok {
		conn.Close()
		delete(clients, userID)
	}
	clientsMu.Unlock()

	broadcastUserStatus()
}

func getClientIDs() []int {
	keys := make([]int, 0, len(clients))
	for key := range clients {
		keys = append(keys, key)
	}
	return keys
}
