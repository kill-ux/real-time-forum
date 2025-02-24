import { socket } from "./app.js";
import { userInfo } from "./auth.js";
import Utils from "./utils.js";

export class ChatManager {
    constructor(receiverUser) {
        this.receiverUser = receiverUser;
        this.messageInput = document.querySelector('.chat-input');
        this.chatMessages = document.querySelector('.chat-body');
        this.setupEventListeners();
        this.loadMessages();
    }

    setupEventListeners() {
        this.messageInput.addEventListener('keypress', e => {
            if (e.shiftKey && e.key === 'Enter') {
                return;
            }
            if (e.key === 'Enter') {
                const value = this.messageInput.value.trim();
                this.messageInput.value = ""

                if (value) {
                    socket.sendMessage({
                        type: "new_message", message: {
                            content: value,
                            receiver_id: this.receiverUser.id,
                            created_at: new Date().getTime()
                        }
                    });
                }

            }
        });

        this.chatMessages.addEventListener('scroll', Utils.opThrottle(() => {
            if (this.chatMessages.scrollTop <= 50) {
                this.loadMessages();
            }
        }, 250));
    }

    async loadMessages() {
        const oldestMessage = this.chatMessages.firstElementChild;
        const before = +oldestMessage?.dataset.timestamp || new Date().getTime();
        const response = await fetch(`/messages`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                receiver_id: this.receiverUser.id,
                before
            })
        });
        const messages = await response.json();
        let oldHeight = this.chatMessages.scrollHeight;
        messages?.forEach((msg) => {
                this.addMessage(msg, true)
        });

        this.chatMessages.scrollTop = this.chatMessages.scrollHeight - oldHeight;
    }

    addMessage(msg, prepend = false, scrollTo) {
        const user = msg.sender_id === this.receiverUser.id ? this.receiverUser : userInfo;
        const element = document.createElement('div');
        element.dataset.timestamp = msg.created_at;
        element.className = `chat-message`;
        element.innerHTML = /*html */`
                    <img src="/assets/images/pics/${Utils.sanitizeHTML(user.image)}" alt="profile-image" class="profile-image">
                    <div class="message-info">
                        <p class="name">${Utils.sanitizeHTML(user.firstname)} ${Utils.sanitizeHTML(user.lastname)}<sub>@${Utils.sanitizeHTML(user.nickname)}</sub></p>
                        <p class="created_at">${new Date(Utils.sanitizeHTML(new Date(msg.created_at))).toLocaleTimeString()}</p>
                        <pre class="message-content">${Utils.sanitizeHTML(msg.content)}</pre>
                    </div>
        `;

        prepend ? this.chatMessages.prepend(element) : this.chatMessages.append(element);
        if (scrollTo) {
            element.scrollIntoView({ behavior: 'smooth' });
        }

    }
}