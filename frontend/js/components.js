import { container } from "./app.js";
import { renderAuthForms } from "./auth.js";
import { renderUsers } from "./chat.js";
import { renderComment } from "./comments.js";
import { renderPost } from "./posts.js";
import Utils from "./utils.js";

/**
 * Renders the main home page layout with sidebar, header, and posts section.
 */
export const renderHome = () => {
    container.innerHTML = /*html */`
        <div class="home">
            <div class="sidebar">
                <div class="users"></div>
            </div>
            <div class="main-content">
                <div class="header">
                    <h1>Real Time Forum</h1>
                    <button class="logout-btn">Logout</button>
                </div>
                <div class="add-post">
                    <form class="addPostForm">
                        <textarea name="content" placeholder="What's on your mind?" required></textarea>
                        <input type="file" name="image" accept="image/*" />
                        <button type="submit">Post</button>
                    </form>
                </div>
                <div class="posts" id="posts"></div>
            </div>
        </div>
    `;
};

/**
 * Renders the list of online users in the sidebar.
 * @param {Array} users - Array of user objects.
 */
export const renderUsers = (users) => {
    const usersAside = document.querySelector(".users")
    usersAside.innerHTML = "<h3>Online Users</h3>"
    users.forEach(user => {
        const userDiv = document.createElement("div")
        userDiv.className = "user-item"
        userDiv.innerHTML = `
            <span>${user.nickname}</span>
            <button class="chat-btn" data-user-id="${user.id}">Chat</button>
        `
        usersAside.appendChild(userDiv)
    })
    document.querySelectorAll(".chat-btn").forEach(btn => {
        btn.addEventListener("click", (e) => {
            const userId = e.target.dataset.userId
            const user = users.find(u => u.id == userId)
            if (user) {
                Utils.openChat(user)
            }
        })
    })
}
