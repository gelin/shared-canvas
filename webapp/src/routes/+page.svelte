<script lang="ts">
    import Canvas from "./Canvas.svelte";
    import ConnectionIndicator from "./ConnectionIndicator.svelte";
    import Palette from "./Palette.svelte";
    import { DEFAULT_TOOL, type PaletteChangeEvent } from "./Palette";

    let tool = $state(DEFAULT_TOOL);
</script>

<section class="card info">
    <nav>
        <p>Draw with the mouse on the canvas below.</p>
        <ConnectionIndicator/>
    </nav>
</section>

<section class="card main">
    <Canvas {tool}/>
    <Palette onPaletteChange={(e: PaletteChangeEvent) => {
        tool = e.tool;
    }}/>
</section>

<style>
    .info {
        font-size: 80%;
        padding: 0 1rem;
    }
    .info nav {
        display: flex;
        flex-direction: row;
        justify-content: space-between;
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
