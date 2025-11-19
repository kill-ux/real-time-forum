/**
 * WebSocket client class for handling real-time communication.
 */
export default class WebWorkerClient {
    /**
     * Initializes the WebSocket client.
     */
    constructor() {
        this.ws = null;
        this.worker = null;
    }

    /**
     * Opens the WebSocket connection using a web worker.
     */
    open() {
        if (this.worker) {
            this.worker.terminate();
        }
        this.worker = new Worker('./worker.js');
        this.worker.postMessage({ type: 'connect', userId: Utils.userId });
        this.worker.onmessage = (e) => {
            const data = e.data;
            switch (data.type) {
                case 'users':
                    renderUsers(data.users);
                    break;
                case 'message':
                    renderMessage(data.message);
                    break;
                case 'posts':
                    getposts();
                    break;
                case 'comments':
                    Utils.getComments(data.postId);
                    break;
            }
        };
    }

    /**
     * Closes the WebSocket connection.
     */
    close() {
        if (this.worker) {
            this.worker.postMessage({ type: 'disconnect' });
            this.worker.terminate();
            this.worker = null;
        }
    }
}
