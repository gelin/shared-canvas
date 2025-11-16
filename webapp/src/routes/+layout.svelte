<script lang="ts">
  import '../app.css';
</script>

<main>
  <h1>Shared Canvas</h1>
  <nav class="card">
    <a href="/" data-sveltekit-preload-data>Home</a>
    <a href="/health" data-sveltekit-preload-data>Health</a>
    <a href="/time" data-sveltekit-preload-data>Time</a>
  </nav>

  <section class="card">
    <slot />
  </section>
</main>
