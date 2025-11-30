<script lang="ts">
    // Home page - no special script needed
    import Canvas from "./Canvas.svelte";
    import ConnectionIndicator from "./ConnectionIndicator.svelte";
</script>

<section class="card info">
    <nav>
        <p>Draw with the mouse on the canvas below.</p>
        <ConnectionIndicator/>
    </nav>
</section>

<section class="card main">
    <Canvas/>
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
    }
</style>
