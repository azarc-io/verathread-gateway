import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import federation from '@originjs/vite-plugin-federation'


// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    // process.env.MODE !== 'production' ? react({
    //   jsxRuntime: 'classic',
    // }) : react(),
    react(),
    federation({
      name: 'app',
      remotes: {
        dummy: "dummy.js",
      },
      shared: ["react", "react-dom", "react-router-dom"],
    })
  ],
  build: {
    sourcemap: true,
    manifest: true,
    target: 'esnext',
  }
})
