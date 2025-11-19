import { container } from "./app.js";
import { renderComment } from "./comments.js";
import Utils from "./utils.js";

/**
 * Renders the HTML structure for a single post.
 * @param {Object} post - The post object containing author, content, image, etc.
 * @returns {string} - The HTML string for the post.
 */
const renderPost = (post) => {
    return /*html*/`
        <div class="post" id="${post.id}">
            <div class="post-header">
                <span class="post-author">${post.author}</span>
                <span class="post-time">${new Date(post.created_at).toLocaleString()}</span>
            </div>
            <div class="post-content">${post.content}</div>
            ${post.image ? `<img src="${post.image}" alt="Post image" class="post-image" />` : ''}
            <div class="post-actions">
                <button class="like-btn" data-post-id="${post.id}">👍 ${post.likes || 0}</button>
                <button class="comment-btn" data-post-id="${post.id}">💬 ${post.comments || 0}</button>
            </div>
            <div class="comments-section" id="comments-${post.id}" style="display: none;">
                <div class="add-comment">
                    <input type="text" placeholder="Add a comment..." class="comment-input" data-post-id="${post.id}" />
                    <button class="submit-comment-btn" data-post-id="${post.id}">Comment</button>
                </div>
                <div class="comments-list" id="comments-list-${post.id}"></div>
            </div>
        </div>
    `;
};

/**
 * Adds event listeners to a post element for like, comment toggle, and submit comment actions.
 * @param {Object} post - The post object.
 * @param {HTMLElement} postElm - The post DOM element.
 */
const addPostEvents = (post, postElm) => {
    postElm.querySelector(`.like`).addEventListener("click", (e) => { Utils.like(+postElm.id, 1, "post_id", e.target.parentElement) })
    postElm.querySelector(`.comment-btn`).addEventListener("click", (e) => {
        const commentsSection = document.getElementById(`comments-${post.id}`);
        commentsSection.style.display = commentsSection.style.display === 'none' ? 'block' : 'none';
    });
    postElm.querySelector(`.submit-comment-btn`).addEventListener("click", (e) => {
        const input = postElm.querySelector(`.comment-input`);
        const content = input.value.trim();
        if (content) {
            Utils.addComment(+postElm.id, content);
            input.value = '';
        }
    });
};

/**
 * Handles the submission of a new post form.
 * @param {Event} e - The form submit event.
 */
export const addPost = async (e) => {
    e.preventDefault();
    const formData = new FormData(e.target);
    const content = formData.get('content').trim();
    if (!content) return;

    const response = await fetch('/posts', {
        method: 'POST',
        body: formData
    });

    if (response.ok) {
        e.target.reset();
        getposts();
    } else {
        Utils.notice('Failed to add post');
    }
};

/**
 * Fetches and renders all posts from the server.
 */
export const getposts = async () => {
    const response = await fetch('/posts');
    const posts = await response.json();
    const postsContainer = document.getElementById('posts');
    postsContainer.innerHTML = '';
    posts.forEach(post => {
        const postElm = document.createElement('div');
        postElm.innerHTML = renderPost(post);
        postsContainer.appendChild(postElm);
        addPostEvents(post, postElm);
        // Load comments for each post
        Utils.getComments(post.id);
    });
};
