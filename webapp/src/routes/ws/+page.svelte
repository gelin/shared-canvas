<script lang="ts">
  import { onDestroy } from 'svelte';
  import { wsClient, type WSMessage } from '$lib/ws';
  import { get } from 'svelte/store';

  let outgoing: string = JSON.stringify({ type: 'ping', payload: { hello: 'world' } }, null, 2);

  const status = wsClient.status;
  const messages = wsClient.messages;

  function connect() {
    wsClient.connect();
  }
  function disconnect() {
    wsClient.disconnect();
  }
  function clear() {
    wsClient.clear();
  }
  function send() {
    try {
      const obj = JSON.parse(outgoing) as WSMessage | Record<string, unknown>;
      const ok = wsClient.send(obj);
      if (!ok) alert('WebSocket not connected');
    } catch (e) {
      alert('Invalid JSON');
    }
  }

  onDestroy(() => {
    // keep connection unless navigating away; do nothing by default
  });
</script>

<h2>WebSocket Demo</h2>
<p>Status: <strong>{$status}</strong></p>
<div style="display:flex; gap:.5rem; flex-wrap: wrap; margin-bottom: .75rem;">
  <button on:click={connect}>Connect</button>
  <button on:click={disconnect}>Disconnect</button>
  <button on:click={clear}>Clear log</button>
  <button on:click={send}>Send JSON</button>
</div>

<div style="display:grid; grid-template-columns: 1fr 1fr; gap: 1rem;">
  <div>
    <h3>Outgoing JSON</h3>
    <textarea bind:value={outgoing} rows={12} style="width:100%; font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New', monospace;"></textarea>
  </div>
  <div>
    <h3>Incoming Messages</h3>
    {#if $messages.length === 0}
      <p>No messages yet.</p>
    {:else}
      {#each $messages as m, i}
        <pre>{JSON.stringify(m, null, 2)}</pre>
      {/each}
    {/if}
  </div>
</div>
