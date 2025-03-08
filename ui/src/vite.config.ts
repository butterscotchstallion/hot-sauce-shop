import {defineConfig, searchForWorkspaceRoot} from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from "@tailwindcss/vite";

// https://vite.dev/config/
export default defineConfig({
    server: {
        fs: {
            allow: [
                // search up for workspace root
                searchForWorkspaceRoot(process.cwd()),
                // your custom rules
                'E:/projects/hot-sauce-shop/ui/node_modules/.pnpm/primeicons@7.0.0/node_modules/primeicons/fonts/primeicons.ttf',
                "E:/projects/hot-sauce-shop/ui/node_modules/.pnpm/primeicons@7.0.0/node_modules/primeicons/fonts/primeicons.woff",
                "E:/projects/hot-sauce-shop/ui/node_modules/.pnpm/primeicons@7.0.0/node_modules/primeicons/fonts/primeicons.woff",
                "E:/projects/hot-sauce-shop/ui/node_modules/.pnpm/primeicons@7.0.0/node_modules/primeicons/fonts/primeicons.woff2"
            ],
        },
    },
    plugins: [react(), tailwindcss()],
})
