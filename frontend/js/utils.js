/**
 * Utility class for common frontend operations.
 */
export default class Utils {
    /**
     * Gets the current user ID from local storage.
     * @returns {string} - The user ID.
     */
    static get userId() {
        return localStorage.getItem('user_id');
    }

    /**
     * Displays a notification message to the user.
     * @param {string} message - The message to display.
     */
    static notice(message) {
        // Implementation for showing notices
    }

    /**
     * Sends a like request for a post or comment.
     * @param {number} id - The ID of the item to like.
     * @param {number} type - The type of like (1 for like).
     * @param {string} column - The column name (post_id or comment_id).
     * @param {HTMLElement} element - The element to update.
     */
    static like(id, type, column, element) {
        // Implementation for liking
    }

    /**
     * Adds a comment to a post.
     * @param {number} postId - The ID of the post.
     * @param {string} content - The comment content.
     */
    static addComment(postId, content) {
        // Implementation for adding comment
    }

    /**
     * Fetches comments for a post.
     * @param {number} postId - The ID of the post.
     */
    static getComments(postId) {
        // Implementation for getting comments
    }

    /**
     * Sends a message to another user.
     * @param {number} userId - The ID of the recipient.
     * @param {string} message - The message content.
     */
    static sendMessage(userId, message) {
        // Implementation for sending message
    }

    /**
     * Opens a chat with a user.
     * @param {Object} user - The user object.
     */
    static openChat(user) {
        // Implementation for opening chat
    }
}
