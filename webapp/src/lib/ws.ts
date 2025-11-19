import { writable, type Readable } from 'svelte/store';

export type WSStatus = 'disconnected' | 'connecting' | 'connected' | 'error';
export type WSMessage = { type: string; payload?: unknown; time?: string };

function wsURL(): string {
  const { protocol, host } = window.location;
  const wsProto = protocol === 'https:' ? 'wss:' : 'ws:';
  return `${wsProto}//${host}/ws`;
}

class WSClient {
  private socket: WebSocket | null = null;
  private statusStore = writable<WSStatus>('disconnected');
  private messagesStore = writable<WSMessage[]>([]);

  get status(): Readable<WSStatus> { return this.statusStore; }
  get messages(): Readable<WSMessage[]> { return this.messagesStore; }

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
        const data = JSON.parse(ev.data) as WSMessage;
        this.messagesStore.update((arr) => [...arr, data]);
      } catch (e) {
        // Non-JSON messages: wrap as text
        this.messagesStore.update((arr) => [...arr, { type: 'text', payload: String(ev.data) }]);
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

  send(obj: WSMessage | Record<string, unknown>): boolean {
    if (!this.socket || this.socket.readyState !== WebSocket.OPEN) return false;
    this.socket.send(JSON.stringify(obj));
    return true;
  }

  clear(): void {
    this.messagesStore.set([]);
  }
}

export const wsClient = new WSClient();
