// shared-worker.js

let socket = null;
const ports = new Set(); // Store all connected tabs

// Handle new connections from tabs
self.onconnect = (event) => {
    const port = event.ports[0];
    ports.add(port);

    // Listen for messages from the tab
    port.onmessage = (event) => {
        const { type, payload } = event.data;

        if (type === 'connect') {
            // Initialize WebSocket connection if it doesn't exist
            if (!socket) {
                socket = new WebSocket('/ws');
                socket.onopen = () => {
                    socket.send(JSON.stringify({ type: 'users' }));
                }

                // Handle incoming messages from the server
                socket.onmessage = (event) => {
                    message = JSON.parse(event.data);
                    // Broadcast the message to all connected tabs
                    ports.forEach((p) => p.postMessage(message));
                };

                // Handle WebSocket errors
                socket.onerror = (error) => {
                    console.error('WebSocket error:', error);
                };

                // Handle WebSocket closure
                socket.onclose = () => {
                    console.log('WebSocket connection closed');
                    socket = null; // Reset the socket
                };
            } else {

                socket.send(JSON.stringify({ type: 'users' }));
            }
        } else if (type === 'send') {
            // Send a message to the server
            if (socket && socket.readyState === WebSocket.OPEN) {
                socket.send(JSON.stringify(payload));
            } else {
                console.error('WebSocket is not open');
            }
        } else if (type === 'read') {
            ports.forEach((p) => p.postMessage(payload));
        } else if (type == "close") {
            ports.clear();
            if (socket) {
                socket.close();
                socket = null;
            }
        }
    };
};
