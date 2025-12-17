<script lang="ts">
    import { wsClient, type WSMessage, type WSUserMessage } from "$lib/ws";
    import { onMount } from "svelte";

    const socket = wsClient;
    let userCount = $state(0);

    onMount(() => {
        socket.connect();
        socket.subscribe((message: WSMessage) => {
            if (message.method === 'user') {
                userCount = message.params.count;
            }
        });
    });
</script>

<p id="user-count" class="user-count" title="{userCount} users online">
    {userCount}
</p>

<style>
    .user-count {
        font-size: 1.4rem;
        padding: 0;
        margin: auto 0;
    }
</style>
