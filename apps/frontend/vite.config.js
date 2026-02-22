import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import { fileURLToPath } from "node:url";

const devApiTarget = process.env.VITE_DEV_API_TARGET || "http://localhost:8080";

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  server: {
    port: 5173,
    strictPort: true,
    proxy: {
      "/api": {
        target: devApiTarget,
        changeOrigin: true,
      },
      "/auth": {
        target: devApiTarget,
        changeOrigin: true,
      },
    },
  },
});
