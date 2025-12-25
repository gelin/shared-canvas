<script lang="ts">
    import Canvas from "./Canvas.svelte";
    import ConnectionIndicator from "./ConnectionIndicator.svelte";
    import Palette from "./Palette.svelte";
    import UserCountIndicator from "./UserCountIndicator.svelte";

    let canvas: Canvas | null = null;
</script>

<section class="card info">
    <nav>
        <div class="left">
            <p>Draw with the mouse on the canvas below</p>
        </div>
        <div class="right">
            <button class="download" title="Download" onclick={() => canvas?.download()}><span class="material-symbols-outlined">download</span></button>
            <UserCountIndicator/>
            <ConnectionIndicator/>
        </div>
    </nav>
</section>

<section class="card main">
    <Canvas bind:this={canvas}/>
    <Palette/>
</section>

<style>
    .info {
        font-size: 80%;
        padding: clamp(0.1rem, 3vw, 0.5rem) clamp(0.1rem, 6vw, 1rem);
    }
    .info p {
        margin: auto 0;
    }
    .info nav, .info nav div {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
    }
    .info nav .right {
        gap: 0.5rem;
    }

    .download {
        font-variation-settings: 'FILL' 0, 'wght' 600, 'GRAD' 0, 'opsz' 32;
        background: none;
        border: none;
        color: inherit;
        padding: 0;
        margin: auto 0;
        cursor: pointer;
    }
    .material-symbols-outlined {
        font-size: 1.4rem;
    }

    :global(.user-count) {
        width: 2rem;
    }

    .main {
        margin-top: 1rem;
        display: flex;
        flex-direction: row;
        gap: clamp(0.5rem, 6vw, 1rem);
    }

    @media (max-width: 800px) {
        .main {
            flex-direction: column;
        }
    }
</style>
