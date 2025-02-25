import { addComment, getComments, renderComment } from "./comments.js";
import { main } from "./components.js";
import Utils from "./utils.js";

// js/posts.js
const renderPost = (post) => {
    return /*html*/`
            <div class="post" id="${post.id}" data-before="${post.created_at}">
            
                <div class="user-profile">
                    <img src="/assets/images/pics/${Utils.sanitizeHTML(post.user.image)}" alt="profile-image" class="profile-image">
                    <div class="profile-info">
                        <p class="name">${Utils.sanitizeHTML(post.user.firstname) + " " + Utils.sanitizeHTML(post.user.lastname)}</p>
                        <p class="role">${Utils.sanitizeHTML(post.user.nickname)}</p>
                    </div>

                    <p class="date">${Utils.sanitizeHTML(new Date(post.created_at * 1000).toLocaleString())}</p>
                </div>

                <div class="post-info">
                    <div class="categories">
                        ${post.categories.split(", ").map(cat => `<span class="category">${Utils.sanitizeHTML(cat)}</span>`).join('')}
                    </div>
                    <h2>${Utils.sanitizeHTML(post.title)}</h2>
                    <p>${Utils.sanitizeHTML(post.content)}</p>
                    
                    ${post.image ? `<img src="/assets/images/posts/${Utils.sanitizeHTML(post.image)}" alt="" class="post-image">` : ''}
                </div>

                <div class="post-footer">

                    <img src="/assets/images/svgs/thumbs-up.svg" alt="" class="like ${post.like === 1 ? 'blue' : ''}">
                    <label class="likes-count">${post.likes}</label>
                    <img src="/assets/images/svgs/thumbs-down.svg" alt="" class="dislike ${post.like === -1 ? 'red' : ''}">
                    <label class="dislikes-count">${post.dislikes}</label>
                    <img src="/assets/images/svgs/comment.svg" class="toggle-comments">
                </div>

            <div class="comments-container">
                <form class="add-comment-form" onsubmit='return false'>
                <input type="text" class="comment-input" data-postId="${post.id}" placeholder="add a comment..." required><button class="add-comment-btn">add</button>
                </form>
                <div class="all-comments">

                </div>
            </div>
            </div>
        `;
}

export const addPost = async (e) => {
    e.preventDefault();

    const formData = new FormData(e.target);
    const title = formData.get("title");
    const content = formData.get("content");
    const categories = formData.getAll("categories");
    const image = formData.get("image");
    if (!title || title.length < 3 || title.length > 100) {
        Utils.notice("Title must be between 3 and 100 characters.");
        return;
    }

    if (!content || content.length < 10 || content.length > 2000) {
        Utils.notice("Content must be between 10 and 2000 characters.");
        return;
    }
    if (image.size > 0 && image.name !== "") {
        if (image.size > 1024 * 1024) { // 1 MB
            Utils.notice("Image size must be less than 1 MB.");
            return;
        }

        const allowedExtensions = ["jpg", "jpeg", "png", "gif", "webp"];
        const imageExtension = image.name.split('.').pop().toLowerCase();

        if (!allowedExtensions.includes(imageExtension)) {
            Utils.notice("Invalid image extension. Allowed extensions are: " + allowedExtensions.join(", "));
            return;
        }
    }

    const allowedCategories = ["tech", "programming", "health", "finance", "food", "science", "memes", "others"];
    for (let category of categories) {
        if (!allowedCategories.includes(category)) {
            Utils.notice("Invalid category: " + category + ". Allowed categories are: " + allowedCategories.join(", "));
            return;
        }
    }


    const response = await fetch('/posts/store', {
        method: 'POST',
        body: formData
    })
    if (response.status === 201) {
        const newPost = await response.json()
        main.insertAdjacentHTML("afterbegin", renderPost(newPost))
        const postElm = document.getElementById(newPost.id)
        addPostEvents(newPost, postElm)
        e.target.reset();
        e.target.parentElement.hidePopover()
        // this.chatMessages.scrollTop = 0
        main.scrollTo({ top: 0, behavior: 'smooth' });
    } else {
        console.log("not added")
    }
}

export const getposts = async (before = +(new Date().getTime() / 1000).toFixed(0)) => {
    const response = await fetch('/posts', {
        method: 'POST',
        body: JSON.stringify({ before })
    })
    const posts = await response.json()
    posts?.forEach((post, i) => {
        main.insertAdjacentHTML("beforeend", renderPost(post))
        const postElm = document.getElementById(post.id)
        addPostEvents(post, postElm)

    });
    if (response.ok) {
        console.log("posts are here")
    } else {
        console.log("posts are not here")

    }
}

export const addCommentsElement = async (id, commentSection, prepend = false) => {
    const before = +commentSection.querySelector(".comment:last-child")?.dataset.before || +(new Date().getTime() / 1000).toFixed(0)
    const comments = await getComments(id, before) || []
    comments.forEach(comment => prepend ? commentSection.prepend(renderComment(comment)) : commentSection.append(renderComment(comment)))
}


const addPostEvents = (post, postElm) => {
    postElm.querySelector(`.like`).addEventListener("click", (e) => { Utils.like(+postElm.id, 1, "post_id", e.target.parentElement) })
    postElm.querySelector(`.dislike`).addEventListener("click", (e) => { Utils.like(+postElm.id, -1, "post_id", e.target.parentElement) })
    postElm.querySelector(`.toggle-comments`).addEventListener("click", (e) => { 
        postElm.querySelector(".post-info").style.cursor = "pointer"
        const commentsContainer = postElm.querySelector(".comments-container")
        commentsContainer.style.width = "80%" 
        const allComments = postElm.querySelector(".all-comments")
        addCommentsElement(post.id, allComments)
        commentsContainer.addEventListener("scroll", Utils.opThrottle((e) => {
            if (e.target.scrollTop + e.target.clientHeight >= e.target.scrollHeight-50) {
                addCommentsElement(post.id, allComments)
            }
        },1000))

    })
    postElm.querySelector(".add-comment-form").addEventListener("submit", addComment)

    postElm.querySelector(".post-info").addEventListener("click", async (e) => {
        postElm.querySelector(".comments-container").style.width == "80%" ? postElm.querySelector(".comments-container").style.width = "0" : ""
        postElm.querySelector(".post-info").style.cursor = "default"
    })
    

}

