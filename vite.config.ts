import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { VitePWA } from 'vite-plugin-pwa'
import mkcert from 'vite-plugin-mkcert'
import fs from 'fs';
import path from 'path';

// ЗАМЕНИТЕ 'RepoName' НА НАЗВАНИЕ ВАШЕГО РЕПОЗИТОРИЯ НА GITHUB
const REPO_NAME = '/acid-app/';

export default defineConfig({
    base: REPO_NAME,
    server: {
        port: 3000,
        https: {
            key: fs.readFileSync(path.resolve(__dirname, 'cert.key')),
            cert: fs.readFileSync(path.resolve(__dirname, 'cert.crt')),
        },
        proxy: {
            "/api": {
                target: "http://localhost:8080",
                changeOrigin: true,
                secure: false, // Важно для проксирования на http с https
            },
        },
    },
    plugins: [
        react(),
        mkcert(),
        VitePWA({
            registerType: 'autoUpdate',
            devOptions: {
                enabled: true,
            },
            manifest: {
                name: "Carbonate Calculator",
                short_name: "AcidApp",
                start_url: REPO_NAME, // Важно указать путь репозитория
                display: "standalone",
                background_color: "#ffffff",
                theme_color: "#004976",
                orientation: "portrait-primary",
                icons: [
                    {
                        "src": "logo192.png", // Вам нужно добавить эти картинки в папку public
                        "type": "image/png",
                        "sizes": "192x192"
                    },
                    {
                        "src": "logo512.png",
                        "type": "image/png",
                        "sizes": "512x512"
                    }
                ]
            }
        })
    ],
})