import { container, socket } from "./app.js"
import { userInfo } from "./auth.js"
import { ChatManager } from "./chat.js"
import { getposts } from "./posts.js"
import Utils from "./utils.js"
export let main

export const renderHome = () => {
    container.innerHTML = /*html */`
    <div class="empty-aside1"></div>
    <div class="nav-aside">
        <nav>
            <!-- profile -->
            <input type="checkbox" id="nav-toggle">
            <div class="nav-profile">
                <img src="/assets/images/pics/${userInfo.image}" alt="profile-image" class="profile-image">

                <div class="profile-info">
                    <p class="name nav-txt">${userInfo.firstname} ${userInfo.lastname}</p>
                    <p class="role nav-txt">${userInfo.nickname}</p>
                </div>
            </div>

            <ul class="nav-items">
                <li onclick="this.classList.toggle('active');this.nextElementSibling.classList.remove('active');document.querySelector('.profile-main').style.display ='none';document.querySelector('.chat-main').style.display ='none';document.querySelector('.home-main').style.display ='grid'"><img src="/assets/images/svgs/home.svg" alt="home-image"><span class='nav-txt'>Home</span></li>
                <li onclick="this.classList.toggle('active');this.previousElementSibling.classList.remove('active');document.querySelector('.profile-main').style.display ='grid';document.querySelector('.chat-main').style.display ='none';document.querySelector('.home-main').style.display ='none'" ><img src="/assets/images/svgs/user.svg" alt="profile-image"><span class='nav-txt'>Profile</span></li>
            </ul>
            <li class="add-post" popovertarget="add-post-form"><button class="add-post-popover-btn"
                    popovertarget="add-post-form"></button><img src="/assets/images/svgs/feather.svg"
                    alt="home-image"><span class='nav-txt'>Post</span></li>
                    <li class="users" >
                    <div class="gradient top"></div>
                    <div class="gradient bottom"></div>
                     </li>
            <li class="logout-btn"><img src="/assets/images/svgs/log-out.svg" alt="home-image"><span class='nav-txt'>logout</span>
            </li>
        </nav>
        <div class="add-post-form" id="add-post-form" popover>
        <button class="close-btn" popovertarget="add-post-form" popovertargetaction="hide">
            <img src="/assets/images/svgs/x.svg" alt="close">
        </button>            
      <form class="addPostForm" enctype="multipart/form-data">
                <label for="">Title :</label>
                <input type="text" name="title" id="">

                <label for="">Content :</label>
                <input type="text" name="content" id="">

                <div class="categories">
                    <label for="tech">Tech</label>
                    <input type="checkbox" name="categories" id="tech" value="tech">

                    <label for="programming">Programming</label>
                    <input type="checkbox" name="categories" id="programming" value="programming">

                    <label for="health">Health</label>
                    <input type="checkbox" name="categories" id="health" value="health">

                    <label for="finance">Finance</label>
                    <input type="checkbox" name="categories" id="finance" value="finance">

                    <label for="food">Food</label>
                    <input type="checkbox" name="categories" id="food" value="food">

                    <label for="science">Science</label>
                    <input type="checkbox" name="categories" id="science" value="science">

                    <label for="memes">Memes</label>
                    <input type="checkbox" name="categories" id="memes" value="memes">

                    <label for="others">Others</label>
                    <input type="checkbox" name="categories" id="others" value="others">
                </div>

                <label for="">Image (optional):</label>
                <input type="file" name="image" accept="image/*" id="">

                <input type="submit" value="Add Post">
            </form>
        </div>
    </div>
    <main class="home-main"> </main>
    <main class="chat-main"> </main>
    <main class="profile-main">
    
    <div class="profile-card">
    <div class="imgContainer">
      <img src="/assets/images/pics/${userInfo.image}" alt="profile card">
    </div>

    <div class="profile-card-cnt">
      <div class="profile-card-name">${userInfo.firstname} ${userInfo.lastname}</div>
      <div class="profile-card-txt">${userInfo.nickname}</div>


      <div class="profile-card-inf">
        <div class="profile-card-inf-item">
          <div class="profile-card-inf-title"> ${new Date(userInfo.created_at*1000).toLocaleDateString()}</div>
          <div class="profile-card-txt">joining date</div>
        </div>

        <div class="profile-card-inf-item">
          <div class="profile-card-inf-title"> ${userInfo.age}</div>
          <div class="profile-card-txt">age</div>
        </div>

        <div class="profile-card-inf-item">
          <div class="profile-card-inf-title"> ${userInfo.gender}</div>
          <div class="profile-card-txt">gender</div>
        </div>

        <div class="profile-card-inf-item">
          <div class="profile-card-inf-title"> ${userInfo.email}</div>
          <div class="profile-card-txt">email</div>
        </div>
      </div>
  </div>

    </main>
    <div class="chat-aside"></div>
    <div class="empty-aside2"></div>
    `
    main = document.querySelector(".home-main")
    main.addEventListener('scroll', Utils.throttle(() => {
        const lastPost = document.querySelector('.post:last-child');
        if (lastPost && lastPost.getBoundingClientRect().top <= window.innerHeight) {
            getposts(+lastPost.dataset.before);
        }
    }, 200));

    document.querySelectorAll(".nav-items li").forEach((item)=>{
        console.log(item)
        item.addEventListener("click",()=>{
            console.log("happen ")
            ChatUI = null
        })
    })
}

export let ChatUI

export const renderUsers = (users) => {
    const usersAside = document.querySelector(".users")
    usersAside.innerHTML = /*html*/ `
            ${users.map(user => /*html*/`
                <div class="user-detail" data-userId="${user.id}">
                    <div class="profile-container user-status ${user.status}">
                    ${user.unread_count > 0 ? `<span class="unread">${user.unread_count > 9 ? "9+" : user.unread_count}</span>` : `<span style="display:none;" class="unread"></span>`}
                    <img  src="/assets/images/pics/${user.image}" alt="home-image" class="profile-image">
                    </div>
                    <span class="span-users nav-txt"> ${user.firstname} ${user.lastname} </span>
                </div>

            `).join("")}
        
    `
    const chatMain = document.querySelector(".chat-main")

    usersAside.childNodes.forEach(element => {
        element.addEventListener("click", () => {

            const userId = element.dataset.userid
            const user = users.find(user => user.id === +userId)
            main.style.display = "none"
            document.querySelector(".profile-main").style.display = "none"
            chatMain.style.display = "grid"
            chatMain.innerHTML = /*html*/ `
                <div class="chat-header">
                    <div class="nav-profile">
                        <img src="/assets/images/pics/${user.image}" alt="profile-image" class="profile-image">

                        <div class="profile-info">
                            <p class="name">${user.firstname} ${user.lastname}</p>
                            <p class="role">${user.nickname}</p>
                        </div>
                    </div>
                </div>
                <div class="chat-body">
                </div>
                </div>
                <div class="chat-footer">
                    <form class="chat-form" onsubmit="return false">
                        <div class="message-container">
                            <textarea id="chat" rows="1" class="message-input chat-input" data-userId="${user.id}" placeholder="Your message..." maxlength="1000"></textarea>
                        </div>
                    </form>
                    <div class="typing u${userId}"><div class="loader"></div><strong>${user.nickname}</strong> is typing...</div>
                </div>
            `
            element.querySelector(".unread").style.display = "none"
            
            ChatUI = new ChatManager(user)

            if (element.dataset.typing === "true"){
                document.querySelector(`.typing`).style.opacity = "1";
            }else{
                document.querySelector(`.typing`).style.opacity = "0";
            }
            socket.markRead(user.id)
        })
    });

}
