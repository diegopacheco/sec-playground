import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

const env = (name: string, fallback: number) => Number(process.env[name] ?? fallback);

export default defineConfig({
  plugins: [react()],
  server: {
    host: "localhost",
    port: env("WEB_PORT", 5175),
    strictPort: true,
    proxy: {
      "/api": { target: `http://127.0.0.1:${env("TRACKER_PORT", 8751)}` },
    },
  },
});
