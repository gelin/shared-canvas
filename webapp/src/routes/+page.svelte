<script lang="ts">
  // Home page - no special script needed
import Canvas from "./Canvas.svelte";
</script>

<section class="card info">
    <p>
        Draw with the mouse on the canvas below.
    </p>
</section>

<section class="card main">
    <Canvas/>
</section>

<style>
    .info {
        font-size: 80%;
        padding: 0 1rem;
    }
    .main {
        margin-top: 1rem;
    }
</style>
