import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// The Go server has no CORS handling, so we proxy /api requests to it in dev.
// Change the target if your server runs on a different host/port.
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      "/api": {
        target: "http://localhost:8000",
        changeOrigin: true,
      },
    },
  },
});
