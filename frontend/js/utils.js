class Utils {
    static(func, wait) {
        let timeout;
        return function (...args) {
            clearTimeout(timeout);
            timeout = setTimeout(() => func.apply(this, args), wait);
        };
    }

    static throttle(func, limit) {
        let inThrottle;
        return (...args) => {
            if (!inThrottle) {
                func(...args);
                inThrottle = true;
                setTimeout(() => inThrottle = false, limit);
            }
        };
    }

    static opThrottle(func, limit) {
        let lastFunc;
        let lastRan;
        return (...args)=> {
            if (!lastRan) {
                func(...args);
                lastRan = Date.now();
            } else {
                clearTimeout(lastFunc);
                lastFunc = setTimeout(() => {
                    if ((Date.now() - lastRan) >= limit) {
                        func(...args);
                        lastRan = Date.now();
                    }
                }, limit - (Date.now() - lastRan));
            }
        };
    }




    static sanitizeHTML(dirtyHTML) {
        const tempDiv = document.createElement('div');
        tempDiv.textContent = dirtyHTML; // Escapes HTML
        return tempDiv.innerHTML;
    }

    static notice(message) {
        const alertMsg = document.createElement("div")
        alertMsg.textContent = message
        alertMsg.className = "alert"
        document.body.append(alertMsg)
        setTimeout(() => {
            alertMsg.remove()
        }, 5000)
    }

    static async like(id, likeOrDislike, commentOrPost, ...postFooters) {
        const response = await fetch('/likes/store', {
            method: 'POST',
            body: JSON.stringify({ p_id: id, name_id: commentOrPost, like: likeOrDislike }),
        })
        if (response.status === 429) {
            Utils.notice("Too Many Requests, slow down!")
            return
        }
        if (response.ok) {
            const likeData = await response.json()
            console.log(postFooters)
            postFooters.forEach((postFooter) => {
                const [likeBtn, likeCount, dislikeBtn, dislikeCount] = postFooter.children
                likeCount.textContent = likeData.likes
                dislikeCount.textContent = likeData.dislikes

                if (likeData.like == 1) {
                    likeBtn.classList.add("blue")
                    dislikeBtn.classList.remove("red")
                } else if (likeData.like === -1) {
                    likeBtn.classList.remove("blue")
                    dislikeBtn.classList.add("red")
                } else {
                    likeBtn.classList.remove("blue")
                    dislikeBtn.classList.remove("red")
                }
            })
        }


    }
}
export default Utils