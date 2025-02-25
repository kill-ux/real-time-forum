import Utils from "./utils.js"

export const renderComment = (comment) => {
    const commentDiv = document.createElement("div")
    commentDiv.className="comment"
    commentDiv.id = comment.id
    commentDiv.dataset.before = comment.created_at
    commentDiv.innerHTML = /*html */ `
    <div class="user-profile" id="${comment.user.id}">
        <img src="/assets/images/pics/${Utils.sanitizeHTML(comment.user.image)}" alt="profile-image" class="profile-image">
        <div class="profile-info">
            <p class="name">${Utils.sanitizeHTML(comment.user.firstname) + " " + Utils.sanitizeHTML(comment.user.lastname)}</p>
            <p class="role">${Utils.sanitizeHTML(comment.user.nickname)}</p>
        </div>

        <p class="date">${Utils.sanitizeHTML(new Date(comment.created_at * 1000).toLocaleString())}</p>
    </div>

    <div class="comment-info">
        <p>${Utils.sanitizeHTML(comment.content)}</p>
    </div>

    <div class="comment-footer">
        <img src="/assets/images/svgs/thumbs-up.svg" alt="" data-id ="${comment.id}"class="like ${comment.like === 1 ? 'blue' : ''}">
        <label class="likes-count">${comment.likes}</label>
        <img src="/assets/images/svgs/thumbs-down.svg" alt="" data-id ="${comment.id}" class="dislike ${comment.like === -1 ? 'red' : ''}">
        <label class="dislikes-count">${comment.dislikes}</label>
    </div>
    `
    const footer = commentDiv.querySelector('.comment-footer');
    commentDiv.querySelector(".like").addEventListener("click", () => Utils.like(+comment.id, 1, 'comment_id', footer))
    commentDiv.querySelector(".dislike").addEventListener("click", () => Utils.like(+comment.id, -1, 'comment_id', footer))
    return commentDiv
}


export const getComments = async (post_id, before) => {
    const response = await fetch('/comments', {
        method: 'POST',
        body: JSON.stringify({ post_id, before })
    })
    const comments = await response.json()
    if (response.ok) {
        return comments
    } else {
        return "something went wrong"
    }
}


export const addComment = async (e) => {
    e.preventDefault()
    const input = e.target.querySelector(".comment-input")
    const content = input.value
    const post_id = +input.dataset.postid
    if (content.length > 2000 || content.length < 1) {
        Utils.notice("comment content length unvalid")
        return
    }
    const response = await fetch('/comments/store', {
        method: 'POST',
        body: JSON.stringify({ content, post_id })
    })
    if (response.status === 201) {
        const newComment = await response.json()
        const commentSection = e.target.parentElement.querySelector(".all-comments")

        commentSection.prepend(renderComment(newComment))

        e.target.reset();
        e.target.parentElement.scrollTo({ top: 0, behavior: 'smooth' });
    } else {
        console.log("not added")
    }
}