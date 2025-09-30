# Real-Time Forum

## Project Overview
This project is a real-time forum application built as part of the Zone01 curriculum. It features user registration and login, post and comment creation, and real-time private messaging using WebSockets. The application is a single-page application (SPA) with a focus on real-time interactions and a responsive user interface.

## Objectives
The goal of this project was to create a modern forum with the following key features:
- **Registration and Login**: Secure user authentication with nickname or email and password.
- **Posts and Comments**: Users can create categorized posts, comment on posts, and view them in a feed.
- **Private Messages**: Real-time chat functionality with online/offline status, message history, and throttled message loading.
- **Real-Time Communication**: Implemented using WebSockets for instant message updates without page refreshes.
- **Single-Page Application**: All page changes are handled dynamically using JavaScript.

## Technologies Used
- **Backend**:
  - **Go**: Handles data processing, WebSocket connections, and server-side logic.
  - **SQLite**: Stores user data, posts, comments, and messages.
- **Frontend**:
  - **JavaScript**: Manages all client-side events and WebSocket communication.
  - **HTML**: Single HTML file for the SPA structure.
  - **CSS**: Stylizes the user interface.
- **No Frontend Frameworks**: Built with vanilla JavaScript, HTML, and CSS as per project requirements.

## Features

### Registration and Login
- **Registration Form**: Users provide:
  - Nickname
  - Age
  - Gender
  - First Name
  - Last Name
  - Email
  - Password
- **Login**: Users can log in using either their nickname or email with their password.
- **Logout**: Available from any page in the forum.
- **Security**: Passwords are hashed using bcrypt.

### Posts and Comments
- **Post Creation**: Users can create posts with categories, similar to the previous forum project.
- **Commenting**: Users can comment on posts, with comments visible only when a post is selected.
- **Feed Display**: Posts are displayed in a feed format for easy browsing.

### Private Messages
- **Chat Interface**:
  - Displays online/offline users, sorted by the last message sent (or alphabetically for new users).
  - Always visible section for selecting chat recipients.
- **Message History**:
  - Loads the last 10 messages when a user is selected.
  - Supports loading 10 additional messages on scroll-up, using throttling/debouncing to prevent scroll event spam.
- **Message Format**:
  - Includes the sender's username.
  - Displays the date and time the message was sent.
- **Real-Time Messaging**: Messages are sent and received instantly via WebSockets.

### Technical Details
- **Single-Page Application**: All navigation and page updates are handled dynamically in JavaScript without reloading the HTML page.
- **WebSockets**:
  - **Backend**: Go with Gorilla WebSocket handles real-time communication.
  - **Frontend**: JavaScript manages client-side WebSocket connections for real-time updates.
- **Database**: SQLite stores all persistent data (users, posts, comments, and messages).
- **Concurrency**: Utilizes Go routines and channels for efficient backend processing.
- **Frontend Optimization**: Throttling and debouncing techniques are applied to optimize scroll events for message loading.

## Learning Outcomes
This project enhanced my understanding of:
- **Web Development**:
  - HTML, CSS, and JavaScript for building a responsive SPA.
  - HTTP, sessions, cookies, and DOM manipulation.
- **Backend Development**:
  - Go programming with routines and channels.
  - WebSocket implementation for real-time communication.
- **Database Management**:
  - SQL and SQLite for data storage and manipulation.
- **Real-Time Systems**:
  - Implementing WebSockets in both Go and JavaScript for seamless real-time interactions.

## Challenges and Solutions
- **Challenge**: Implementing real-time messaging without overwhelming the server or client.
  - **Solution**: Used WebSockets with efficient message handling and throttling/debouncing for scroll events.
- **Challenge**: Managing a single-page application without frameworks.
  - **Solution**: Developed a modular JavaScript structure to handle dynamic page updates and state management.
- **Challenge**: Ensuring secure user authentication.
  - **Solution**: Used bcrypt for password hashing and UUID for unique user identification.

## How to Run
1. **Clone the Repository**:
   ```bash
   git clone https://learn.zone01oujda.ma/git/muboutoub/real-time-forum
   ```
2. **Set Up Dependencies**:
   - Ensure Go is installed and the required packages (Gorilla WebSocket, SQLite3, bcrypt, UUID) are available.
   - No additional frontend dependencies are required.
3. **Run the Backend**:
   ```bash
   go run main.go
   ```
4. **Access the Forum**:
   - http://localhost:8080

## Acknowledgments
This project was completed as part of the Zone01 curriculum, building on the skills learned from the previous forum project. Special thanks to the Zone01 team for providing the opportunity to work on this challenging and rewarding project.