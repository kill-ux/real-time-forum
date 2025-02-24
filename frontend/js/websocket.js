import { ChatUI, renderUsers } from "./components.js";
class WebWorkerClient {
    constructor() {
    }
    open = () => {
        this.worker = new SharedWorker("/js/worker.js")
        this.worker.port.start()
        this.worker.port.postMessage({ type: "connect" })
        this.setupEventListeners();
    }

    setupEventListeners = () => {

        this.worker.port.onmessage = this.handleMessage;

    }

    PreRenderUsers = ({ members, data }) => {
        members = members.map((member) => {
            member.status = data.includes(member.id) ? 'online' : 'offline';
            return member;
        })
        renderUsers(members);
    }

    handleMessage = ({ data }) => {
        switch (data.type) {
            case 'users':
                this.PreRenderUsers(data);
                break;
            case 'status_update':
                console.log(data)
                this.PreRenderUsers(data);
                break;
            case 'read':
                document.querySelector(`.user-detail[data-userId="${data.message.receiver_id}"] .unread`).style.display = "none";
                break
            case 'new_message':
                console.log("data => ", data)
                this.PreRenderUsers(data);
                if (ChatUI?.receiverUser.id === data.message.sender_id || ChatUI?.receiverUser.id === data.message.receiver_id) {
                    ChatUI.addMessage(data.message, false, true);
                    if (data.message.sender_id === ChatUI.receiverUser.id) {
                        this.markRead(data.message.sender_id);
                    }
                }
                break;

        }
    }

    markRead = (receiver_id) => {
        this.sendMessage({ type: 'read', message: { receiver_id } });
        this.worker.port.postMessage({ type: 'read', payload: { type: 'read', message: { receiver_id } } })
    }

    sendMessage = (message) => {
        this.worker.port.postMessage({ type: "send", payload: message });
    }

    getUsers = () => {
        this.sendMessage({ type: 'users' });
    }

    close = () => {
        this.worker.port.postMessage({ type: "close" })
    }


}

export default WebWorkerClient