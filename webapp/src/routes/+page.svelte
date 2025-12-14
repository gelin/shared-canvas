<script lang="ts">
    import Canvas from "./Canvas.svelte";
    import ConnectionIndicator from "./ConnectionIndicator.svelte";
    import Palette from "./Palette.svelte";
    import { DEFAULT_TOOL, type PaletteChangeEvent } from "./Palette";

    let canvas: Canvas;
    let tool = $state(DEFAULT_TOOL);
</script>

<section class="card info">
    <nav>
        <div class="left">
            <p>Draw with the mouse on the canvas below.</p>
        </div>
        <div class="right">
            <button class="download" title="Download" onclick={() => canvas?.download()}><span class="material-symbols-outlined">download</span></button>
            <ConnectionIndicator/>
        </div>
    </nav>
</section>

<section class="card main">
    <Canvas {tool} bind:this={canvas}/>
    <Palette onPaletteChange={(e: PaletteChangeEvent) => {
        tool = e.tool;
    }}/>
</section>

<style>
    .info {
        font-size: 80%;
        padding: 0 1rem;
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
        font-variation-settings: 'FILL' 0, 'wght' 300, 'GRAD' 0, 'opsz' 48;
        background: none;
        border: none;
        color: inherit;
        padding: 0;
        cursor: pointer;
    }

    .main {
        margin-top: 1rem;
        display: flex;
        flex-direction: row;
        gap: 1rem;
    }

    @media (max-width: 800px) {
        .main {
            flex-direction: column;
        }
    }
</style>
