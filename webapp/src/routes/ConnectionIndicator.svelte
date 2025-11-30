<script lang="ts">
    import { wsClient, type WSStatus } from "$lib/ws";

    const status = wsClient.status;
</script>

<p>
{#if $status === 'disconnected'}
    <span class="disconnected" title="Disconnected">⚫</span>
{:else if $status === 'connecting'}
    <span class="connecting" title="Connecting...">🟡</span>
{:else if $status === 'connected'}
    <span class="connected" title="Connected">🟢</span>
{:else if $status === 'error'}
    <span class="error" title="Error">🔴</span>
{/if}
</p>

<style>
    p {
        cursor: default;
    }
</style>
