import { userInfo } from "./auth.js";
import { ChatUI, renderUsers } from "./components.js";
import Utils from "./utils.js";

class WebWorkerClient {
    constructor() {
    }
    open = () => {
        this.worker = new SharedWorker("/js/worker.js")
        this.worker.port.start()
        this.worker.port.postMessage({ type: "connect" })
        this.setupEventListeners();
        this.opacity
        this.timer
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
                this.PreRenderUsers(data);
                break;
            case 'read':
                document.querySelector(`.user-detail[data-userId="${data.message.receiver_id}"] .unread`).style.display = "none";
                break
            case 'error':
                Utils.notice("message was not sent, try again!")
                break
            case 'ping':
                this.worker.port.postMessage({ type: "pong" })
                break
            case 'typing':

                const user = document.querySelector(`.user-detail[data-userid="${data.message.sender_id}"]`)
                if (user) {
                    user.dataset.typing = data.is_typing
                }
                //userInfo.id !== data.message.sender_id && userInfo.id === data.message.receiver_id && 
                if (data.message.sender_id == ChatUI?.receiverUser.id) {
                    if (data.is_typing) {
                        this.opacity = "1";
                    } else {
                        this.opacity = "0";
                    }
                    clearTimeout(this.timer)
                    this.timer = setTimeout(() =>
                        document.querySelector(`.typing.u${data.message.sender_id}`).style.opacity = this.opacity
                        , 300)
                }

                break
            case 'new_message':
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