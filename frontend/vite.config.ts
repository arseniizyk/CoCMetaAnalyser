import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'

// https://vite.dev/config/
export default defineConfig({
    base: '/cocmetaanalyser/', // <-- добавь это
    plugins: [react()],
})
