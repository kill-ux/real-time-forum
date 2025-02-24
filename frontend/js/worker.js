// shared-worker.js

let socket = null;
const ports = new Set(); // Store all connected tabs
const lastPong = new WeakMap(); // Track last pong time per port

// Heartbeat configuration
const PING_INTERVAL = 1000; // Send ping every 5 seconds
const PONG_TIMEOUT = 2000; // Consider port closed if no pong in 10 seconds

// Function to ping all ports and check for timeouts
function pingPorts() {
    const now = Date.now();
    ports.forEach((port) => {
        try {
            port.postMessage({ type: 'ping' });
            const last = lastPong.get(port);
            if (last && now - last > PONG_TIMEOUT) {
                console.log('Port timed out, removing it');
                ports.delete(port);
                // socket.sendMessage({
                //     type: "typing", message: {
                //         receiver_id: this.receiverUser.id
                //     },
                //     is_typing: false
                // });
                if (ports.size === 0 && socket) {
                    socket.close();
                    socket = null;
                    console.log('WebSocket closed due to no active ports');
                }

            }
        } catch (e) {
            console.error('Error pinging port:', e);
            ports.delete(port); // Remove port if messaging fails
            if (ports.size === 0 && socket) {
                socket.close();
                socket = null;
            }
        }
    });
    setTimeout(pingPorts, PING_INTERVAL); // Schedule next ping
}

// Start the ping loop
pingPorts();

// Handle new connections from tabs
self.onconnect = (event) => {
    const port = event.ports[0];
    ports.add(port);
    lastPong.set(port, Date.now()); // Initialize last pong time

    // Listen for messages from the tab
    port.onmessage = (event) => {
        const { type, payload } = event.data;

        if (type === 'pong') {
            lastPong.set(port, Date.now()); // Update last pong time
        } else if (type === 'connect') {
            if (!socket) {
                socket = new WebSocket('/ws');
                socket.onopen = () => {
                    socket.send(JSON.stringify({ type: 'users' }));
                };
                socket.onmessage = (event) => {
                    const message = JSON.parse(event.data);
                    ports.forEach((p) => p.postMessage(message));
                };
                socket.onerror = (error) => {
                    console.error('WebSocket error:', error);
                };
                socket.onclose = () => {
                    console.log('WebSocket connection closed');
                    socket = null;
                };
            } else {
                socket.send(JSON.stringify({ type: 'users' }));
            }
        } else if (type === 'send') {
            if (socket && socket.readyState === WebSocket.OPEN) {
                socket.send(JSON.stringify(payload));
            } else {
                console.error('WebSocket is not open');
            }
        } else if (type === 'read') {
            ports.forEach((p) => p.postMessage(payload));
        } else if (type === 'close') {
            ports.clear();
            if (socket) {
                socket.close();
                socket = null;
            }
        }
    };
};