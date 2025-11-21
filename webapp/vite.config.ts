import { defineConfig } from 'vite';
import { sveltekit } from '@sveltejs/kit/vite';

export default defineConfig({
    plugins: [sveltekit()],
    server: {
        port: 5173,
        proxy: {
            '/socket': {
                target: 'http://localhost:8080',
                changeOrigin: true,
                ws: true,
            },
            '/health': {
                target: 'http://localhost:8080',
                changeOrigin: true,
            },
        },
    },
});
