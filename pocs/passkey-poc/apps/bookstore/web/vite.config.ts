import { defineConfig } from "vite";

const env = (name: string, fallback: number) => Number(process.env[name] ?? fallback);

export default defineConfig({
  server: {
    host: "localhost",
    port: env("BOOKSTORE_WEB_PORT", 5173),
    strictPort: true,
    proxy: {
      "/passkey": { target: `http://127.0.0.1:${env("PASSKEY_PORT", 8090)}`, rewrite: (p) => p.replace(/^\/passkey/, "/api") },
      "/api": { target: `http://127.0.0.1:${env("BOOKSTORE_API_PORT", 4001)}` },
    },
  },
});
