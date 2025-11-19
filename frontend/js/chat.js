import { container } from "./app.js";
import { renderHome } from "./components.js";
import Utils from "./utils.js";

/**
 * Renders the chat interface for a specific user.
 * @param {Object} user - The user object to chat with.
 */
export const renderChat = (user) => {
    container.innerHTML = /*html*/`
    <div class="chat-container">
        <div class="chat-header">
            <button class="back-btn">← Back</button>
            <h2>Chat with ${user.nickname}</h2>
        </div>
        <div class="chat-messages" id="chat-messages"></div>
        <div class="chat-input">
            <input type="text" id="message-input" placeholder="Type a message..." />
            <button id="send-btn">Send</button>
        </div>
    </div>
    `;

    document.querySelector(".back-btn").addEventListener("click", () => {
        renderHome();
    });

    const sendBtn = document.getElementById("send-btn");
    const messageInput = document.getElementById("message-input");

    sendBtn.addEventListener("click", () => {
        const message = messageInput.value.trim();
        if (message) {
            Utils.sendMessage(user.id, message);
            messageInput.value = "";
        }
    });

    messageInput.addEventListener("keypress", (e) => {
        if (e.key === "Enter") {
            sendBtn.click();
        }
    });
};

/**
 * Renders a single message in the chat interface.
 * @param {Object} message - The message object to render.
 */
export const renderMessage = (message) => {
    const chatMessages = document.getElementById("chat-messages");
    if (!chatMessages) return;

    const messageDiv = document.createElement("div");
    messageDiv.className = `message ${message.sender_id === Utils.userId ? "own" : "other"}`;
    messageDiv.innerHTML = `
        <div class="message-content">${message.content}</div>
        <div class="message-time">${new Date(message.created_at).toLocaleTimeString()}</div>
    `;
    chatMessages.appendChild(messageDiv);
    chatMessages.scrollTop = chatMessages.scrollHeight;
};

/**
 * Renders the list of online users in the sidebar.
 * @param {Array} users - Array of user objects.
 */
export const renderUsers = (users) => {
    const usersAside = document.querySelector(".users");
    if (!usersAside) return;

    usersAside.innerHTML = "<h3>Online Users</h3>";
    users.forEach(user => {
        const userDiv = document.createElement("div");
        userDiv.className = "user-item";
        userDiv.innerHTML = `
            <span>${user.nickname}</span>
            <button class="chat-btn" data-user-id="${user.id}">Chat</button>
        `;
        usersAside.appendChild(userDiv);
    });

    document.querySelectorAll(".chat-btn").forEach(btn => {
        btn.addEventListener("click", (e) => {
            const userId = e.target.dataset.userId;
            const user = users.find(u => u.id == userId);
            if (user) {
                renderChat(user);
            }
        });
    });
};
