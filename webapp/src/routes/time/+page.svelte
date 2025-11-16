<script lang="ts">
  import { fetchJSON } from '$lib/api';
  let output: string = 'Click to get server time...';
  let loading = false;
  async function getTime() {
    loading = true;
    output = 'Loading...';
    try {
      const data = await fetchJSON('/api/time');
      output = JSON.stringify(data, null, 2);
    } catch (e: any) {
      output = 'Error: ' + (e?.message ?? e);
    } finally {
      loading = false;
    }
  }
</script>

<h2>Time</h2>
<button on:click={getTime} disabled={loading}>{loading ? 'Loading...' : 'Get /api/time'}</button>
<pre>{output}</pre>
