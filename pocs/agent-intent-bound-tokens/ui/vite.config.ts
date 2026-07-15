import { defineConfig, loadEnv } from "vite"
import react from "@vitejs/plugin-react"

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, ".", "")
  const api = `http://localhost:${env.PORT || "8081"}`

  return {
    plugins: [react()],
    server: {
      port: Number(env.UI_PORT || "5174"),
      proxy: {
        "/api": api,
        "/health": api,
        "/service-index": {
          target: api,
          rewrite: () => "/",
        },
      },
    },
  }
})
