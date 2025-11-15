async function fetchJSON(url) {
  const res = await fetch(url, { headers: { 'Accept': 'application/json' } });
  if (!res.ok) throw new Error(`HTTP ${res.status}`);
  return await res.json();
}

document.addEventListener('DOMContentLoaded', () => {
  const pingBtn = document.getElementById('ping');
  const healthOutput = document.getElementById('healthOutput');
  const timeBtn = document.getElementById('getTime');
  const timeOutput = document.getElementById('timeOutput');

  pingBtn.addEventListener('click', async () => {
    healthOutput.textContent = 'Loading...';
    try {
      const data = await fetchJSON('/api/health');
      healthOutput.textContent = JSON.stringify(data, null, 2);
    } catch (e) {
      healthOutput.textContent = 'Error: ' + e.message;
    }
  });

  timeBtn.addEventListener('click', async () => {
    timeOutput.textContent = 'Loading...';
    try {
      const data = await fetchJSON('/api/time');
      timeOutput.textContent = JSON.stringify(data, null, 2);
    } catch (e) {
      timeOutput.textContent = 'Error: ' + e.message;
    }
  });
});
