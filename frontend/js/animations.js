let switchCtn 
let switchC1 
let switchC2 
let switchBtn
let aContainer 
let bContainer 
let allButtons 

let getButtons = (e) => e.preventDefault();

export const changeForm = (e) => {
    switchCtn.classList.add("is-gx");
    setTimeout(function () {
        switchCtn.classList.remove("is-gx");
    }, 1500);

    switchCtn.classList.toggle("is-txr");

    switchC1.classList.toggle("is-hidden");
    switchC2.classList.toggle("is-hidden");
    aContainer.classList.toggle("is-txl");
    bContainer.classList.toggle("is-txl");
    bContainer.classList.toggle("is-z200");
};

export let mainF = () => {
    switchCtn = document.querySelector("#switch-cnt");
    switchC1 = document.querySelector("#switch-c1");
    switchC2 = document.querySelector("#switch-c2");
    switchBtn = document.querySelectorAll(".switch-btn");
    aContainer = document.querySelector("#a-container");
    bContainer = document.querySelector("#b-container");
    allButtons = document.querySelectorAll(".submit");
    for (var i = 0; i < allButtons.length; i++)
        allButtons[i].addEventListener("click", getButtons);
    for (var i = 0; i < switchBtn.length; i++)
        switchBtn[i].addEventListener("click", changeForm);
};
// window.addEventListener("load", mainF);