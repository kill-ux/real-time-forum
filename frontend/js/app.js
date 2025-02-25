import { mainF } from "./animations.js";
import { checkUserLogin, domLogout, handleLogout, isLoggedin, renderAuthForms } from "./auth.js";
import { renderHome } from "./components.js";
import { addPost, getposts } from "./posts.js";
import Utils from "./utils.js";
import WebWorkerClient from "./websocket.js";


export const container = document.querySelector(".container")
export const socket = new WebWorkerClient()


// Wait for the DOM to load
document.addEventListener("DOMContentLoaded", async () => {



  // Save the original fetch function
  const originalFetch = window.fetch;

  // Override global fetch
  window.fetch = async function (...args) {

    try {
      const response = await originalFetch(...args);
      // Automatically check for 429 status
      if (response.status === 429) {
        // Cancel the request chain
        console.log('429 Rate Limit Exceeded - Request Aborted');
        Utils.notice("Too Many Requests, slow down!")
      } else if (response.status == 401) {
        domLogout()
      }else if (response.status == 500) {
        Utils.notice("something went wrong")
      }

      // For non-429 responses, return original response
      return response;
    } catch (error) {
      console.log(error)
    }
    // Return a dummy response in case of an error
    return new Response(JSON.stringify({ error: "Network error" }), {
      status: 500,
      headers: { "Content-Type": "application/json" }
    });
  };
  const notFound = new Promise((resolve)=>{
    if (window.location.pathname === "/"){
      resolve("")
    }else{
      const closeBtn = document.createElement("button")
      closeBtn.textContent = "continue"
      closeBtn.onclick = ()=>{
        history.replaceState(null, "", "/");
        resolve("")
      }
      container.innerHTML = /*html*/ `
      <div class="Error404"><div class="innerDiv"><h1> 404 not Found </h1></div></div>
      `
      document.querySelector(".innerDiv").append(closeBtn)
    }

  })
  await notFound
  await checkUserLogin()

  if (isLoggedin) {

    renderHome()
    document.querySelector(".logout-btn").addEventListener("click", handleLogout)
    document.querySelector(".addPostForm").addEventListener("submit", addPost)
    getposts()
    socket.open()
  } else {
    renderAuthForms()
    mainF()
  }


});


