import { container } from "./app.js";
import Utils from "./utils.js";

/**
 * Renders a single comment element.
 * @param {Object} comment - The comment object containing author, content, likes, etc.
 * @returns {HTMLElement} - The rendered comment div element.
 */
export const renderComment = (comment) => {
    const commentDiv = document.createElement("div")
    commentDiv.className = "comment"
    commentDiv.innerHTML = /*html*/`
        <div class="comment-header">
            <span class="comment-author">${comment.author}</span>
            <span class="comment-time">${new Date(comment.created_at).toLocaleString()}</span>
        </div>
        <div class="comment-content">${comment.content}</div>
        <div class="comment-actions">
            <button class="like-btn" data-comment-id="${comment.id}">👍 ${comment.likes || 0}</button>
        </div>
    `
    commentDiv.querySelector(".like-btn").addEventListener("click", (e) => {
        Utils.like(+e.target.dataset.commentId, 1, "comment_id", e.target)
    })
    return commentDiv
}
