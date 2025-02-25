import { changeForm, mainF } from "./animations.js";
import { container, socket } from "./app.js";
import { renderHome } from "./components.js";
import { addPost, getposts } from "./posts.js";
import Utils from "./utils.js";

export let isLoggedin = false
export let userInfo

const patterns = {
    nickname: /^[a-zA-Z0-9_-]{3,40}$/,
    email: /^[a-zA-Z0-9._+-]{3,20}@[a-zA-Z0-9.-]{3,20}\.[a-zA-Z]{2,10}$/,
    password: /^.{8,100}$/,
    firstName: /^[A-Za-z]{3,40}$/,
    lastName: /^[A-Za-z]{3,40}$/,
    gender: /^(male|female)$/,
};

const renderTemplate = (html) => {
    const container = document.getElementById("authContainer")
    container.innerHTML = html;
    if (html) {
        container.style.display = "flex"
    } else {
        container.style.display = "none"
    }
}

export const renderAuthForms = () => {
    renderTemplate(/*html */`
        <div class="log">
        <div>
          <div class="auth-container a-container" id="a-container">
            <form class="form" id="a-form" onsubmit="return false">
              <h2 class="form_title title">Create Account</h2>
              <input
                class="form__input"
                type="text"
                id="nickname"
                placeholder="Nickname"
              />
              <input
                class="form__input"
                type="number"
                id="age"
                placeholder="Age"
              />
              <div>
                <label for="male" class="label-male">Male</label>
                <input type="radio" name="gender" id="male" value="male" />
                <label for="female" class="label-male">Female</label>
                <input type="radio" name="gender" id="female" value="female" />
              </div>
              <input
                class="form__input"
                type="text"
                id="firstname"
                placeholder="First Name"
              />
              <input
                class="form__input"
                type="text"
                id="lastname"
                placeholder="Last Name"
              />
              <input
                class="form__input"
                type="email"
                placeholder="Email"
                id="signup-email"
              />
              <input
                class="form__input"
                type="password"
                placeholder="Password"
                id="signup-password"
              />
              <button id="signupBtn" class="form__button button submit">
                Sign Up
              </button>
            </form>
          </div>
          <div class="auth-container b-container" id="b-container">
            <form class="form" id="b-form" onsubmit="return false">
              <h2 class="form_title title">Sign in to Website</h2>
              <input
                class="form__input"
                type="text"
                placeholder="Email / Nickname"
                id="login-email"
                required
              />
              <input
                class="form__input"
                type="password"
                placeholder="Password"
                id="login-password"
                required
              />
              <button class="form__button button submit" id="loginBtn">
                SIGN IN
              </button>
            </form>
          </div>
        </div>
        <div class="switch" id="switch-cnt">
          <div class="switch__container" id="switch-c1">
            <h2 class="switch__title title">Welcome Back !</h2>
            <p class="switch__description description">
              To keep connected with us please login with your personal info
            </p>
            <button class="switch__button button switch-btn">SIGN IN</button>
          </div>
          <div class="switch__container is-hidden" id="switch-c2">
            <h2 class="switch__title title">Hello Friend !</h2>
            <p class="switch__description description">
                Enter your personal details and start journey with us in real time forum.
            </p>
            <button class="switch__button button switch-btn">SIGN UP</button>
          </div>
        </div>
        `)

    document.getElementById("loginBtn").addEventListener("click", handleLogin);
    document.getElementById("signupBtn").addEventListener("click", handleSignup);


}

const validateInput = (id, pattern) => {
    const input = document.getElementById(id)
    const value = input.value.trim();
    if (!value || !pattern.test(value)) {
        input.classList.add("Error")
        return false;
    }
    input.classList.remove("Error")
    return value;
}

const handleLogin = async () => {
    const email = validateInput("login-email", patterns.nickname) ||
        validateInput("login-email", patterns.email);
    const password = validateInput("login-password", patterns.password);
    if (!email || !password) return;

    const response = await fetch('/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
    });

    const data = await response.json();
    Utils.notice(data.message);

    if (response.ok) {
        userInfo = data.user
        renderHome()

        document.querySelector(".logout-btn").addEventListener("click", handleLogout)
        document.querySelector(".addPostForm").addEventListener("submit", addPost)
        renderTemplate()
        localStorage.setItem("user_id", data.user.id)
        getposts()
        socket.open()
    }
}
export const domLogout = ()=>{
  userInfo = undefined
  Utils.notice("goodbye")
  renderAuthForms()
   mainF()
  container.innerHTML = ""
  localStorage.removeItem("user_id")
  socket.close()
  isLoggedin = false
}
export const handleLogout = async () => {
    const response = await fetch('/logout');
    if (response.status === 204) {
      domLogout()
    } else {
        Utils.notice("didn't work")
    }
}


const handleSignup = async (e) => {
    const nickname = validateInput("nickname", patterns.nickname, "*Invalid nickname");
    const ageInput = document.getElementById("age");
    const age = +ageInput.value;
    age >= 10 && age <= 200 ? ageInput.classList.remove("Error") : ageInput.classList.add("Error");
    const gender = document.querySelector("input[name='gender']:checked")?.value;
    const genderInputs = document.querySelectorAll("label");
    patterns.gender.test(gender) ? genderInputs.forEach((radio)=>radio.classList.remove("Error")) : genderInputs.forEach((radio)=>radio.classList.add("Error")) ;
    const firstname = validateInput("firstname", patterns.firstName, "*Invalid first name");
    const lastname = validateInput("lastname", patterns.lastName, "*Invalid last name");
    const email = validateInput("signup-email", patterns.email, "*Invalid email");
    const password = validateInput("signup-password", patterns.password, "*Invalid password");

    if (!nickname || !age || !gender || !firstname || !lastname || !email || !password) return;

    const response = await fetch('/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ nickname, age, gender, firstname, lastname, email, password })
    });

    const data = await response.json();
    Utils.notice(data.message);
    if (response.status === 201) {
      changeForm()
      e.target.parentElement.reset()
    };
}


export const checkUserLogin = async () => {
    const response = await fetch('/check-auth');
    const data = await response.json()
    if (response.ok) {
        userInfo = data
        isLoggedin = true
    } else {
        isLoggedin = false
    }
}