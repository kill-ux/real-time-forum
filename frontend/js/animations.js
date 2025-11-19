/**
 * Prevents the default action of an event.
 * @param {Event} e - The event object.
 */
let getButtons = (e) => e.preventDefault();

/**
 * Handles the form switching animation between login and signup forms.
 * @param {Event} e - The event object.
 */
export const changeForm = (e) => {
    switchCtn.classList.add("is-gx");
    setTimeout(function(){
        switchCtn.classList.remove("is-gx");
    }, 1500)

    switchCtn.classList.toggle("is-txr");
    aContainer.classList.toggle("is-txl");
    bContainer.classList.toggle("is-txl");
    bContainer.classList.toggle("is-z200");
}

/**
 * Initializes the animation components and event listeners for form switching.
 */
export let mainF = () => {
    switchCtn = document.querySelector("#switch-cnt");
    aContainer = document.querySelector("#a-container");
    bContainer = document.querySelector("#b-container");
    allButtons = document.querySelectorAll(".switch-btn");
    allButtons.forEach((btn)=>{
        btn.addEventListener("click", changeForm)
    });
    /**
     * Prevents the default action of an event inside mainF.
     * @param {Event} e - The event object.
     */
    let getButtons = (e) => e.preventDefault();
}
