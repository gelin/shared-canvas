import { defineConfig } from 'vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { SvelteKitPWA } from '@vite-pwa/sveltekit';

export default defineConfig({
    plugins: [
        sveltekit(),
        SvelteKitPWA({
            registerType: 'autoUpdate',
            injectRegister: 'auto',
            manifest: {
                name: 'Shared Canvas',
                short_name: 'SharedCanvas',
                description: 'A shared canvas for drawing together',
                theme_color: 'black',
                background_color: 'black',
                display: 'standalone',
                scope: '/',
                start_url: '/',
                icons: [
                    {
                        src: 'icon-192x192.png',
                        sizes: '192x192',
                        type: 'image/png'
                    },
                    {
                        src: 'icon-512x512.png',
                        sizes: '512x512',
                        type: 'image/png'
                    },
                    {
                        src: 'icon-1024x1024.png',
                        sizes: '1024x1024',
                        type: 'image/png',
                    }
                ],
                screenshots: [
                    {
                        src: 'screenshot-mobile.png',
                        sizes: '1024x1512',
                        type: 'image/png',
                        form_factor: 'narrow',
                        label: 'Shared Canvas on mobile'
                    },
                    {
                        src: 'screenshot-desktop.png',
                        sizes: '2048x1572',
                        type: 'image/png',
                        form_factor: 'wide',
                        label: 'Shared Canvas on desktop'
                    }
                ]
            }
        })
    ],
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
