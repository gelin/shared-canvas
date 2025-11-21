import { writable, type Readable } from 'svelte/store';

export type WSStatus = 'disconnected' | 'connecting' | 'connected' | 'error';
export type WSDrawMessage = { method: 'draw'; params: { x: number; y: number; w: number; h: number; p: string }; };
export type WSCallback = (msg: WSDrawMessage) => void;

function wsURL(): string {
    const { protocol, host } = window.location;
    const wsProto = protocol === 'https:' ? 'wss:' : 'ws:';
    return `${wsProto}//${host}/socket`;
}

class WSClient {
    private socket: WebSocket | null = null;
    private statusStore = writable<WSStatus>('disconnected');
    private subscribers: WSCallback[] = [];

    get status(): Readable<WSStatus> {
        return this.statusStore;
    }

    subscribe(cb: WSCallback): void {
        this.subscribers.push(cb);
    }

    connect(): void {
        if (this.socket && (this.socket.readyState === WebSocket.OPEN || this.socket.readyState === WebSocket.CONNECTING)) return;
        this.statusStore.set('connecting');
        const url = wsURL();
        const ws = new WebSocket(url);
        this.socket = ws;

        ws.onopen = () => {
            this.statusStore.set('connected');
        };

        ws.onmessage = (ev) => {
            try {
                const data = JSON.parse(ev.data) as WSDrawMessage;
                for (const cb of this.subscribers) {
                    cb(data);
                }
            } catch (e) {
                console.error('Failed to process message:', ev.data, e);
            }
        };

        ws.onerror = () => {
            this.statusStore.set('error');
        };

        ws.onclose = () => {
            this.statusStore.set('disconnected');
        };
    }

    disconnect(): void {
        if (this.socket) {
            this.socket.close(1000, 'client closing');
            this.socket = null;
        }
    }

    send(obj: WSDrawMessage | Record<string, unknown>): boolean {
        if (!this.socket || this.socket.readyState !== WebSocket.OPEN) return false;
        this.socket.send(JSON.stringify(obj));
        return true;
    }
}

export const wsClient = new WSClient();
