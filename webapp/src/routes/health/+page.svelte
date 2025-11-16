<script lang="ts">
  import { fetchJSON } from '$lib/api';
  let output: string = 'Click to ping...';
  let loading = false;
  async function ping() {
    loading = true;
    output = 'Loading...';
    try {
      const data = await fetchJSON('/api/health');
      output = JSON.stringify(data, null, 2);
    } catch (e: any) {
      output = 'Error: ' + (e?.message ?? e);
    } finally {
      loading = false;
    }
  }
</script>

<h2>Health</h2>
<button on:click={ping} disabled={loading}>{loading ? 'Pinging...' : 'Ping /api/health'}</button>
<pre>{output}</pre>
