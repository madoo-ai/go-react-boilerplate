import { defineConfig } from "vite"
import react from "@vitejs/plugin-react"
import tailwindcss from "@tailwindcss/vite"
import { fileURLToPath, URL } from "node:url"

export default defineConfig({
  plugins: [
    react(),
    tailwindcss(),
    {
      name: "app-config-dev",
      apply: "serve",
      transformIndexHtml(html) {
        const locale = process.env.APP_LOCALE === "en" ? "en" : "id"
        const name = process.env.APP_NAME ?? "Go React Boilerplate"
        return html
          .replace('"__APP_LOCALE__"', `"${locale}"`)
          .replace('"__APP_NAME__"', `"${name.replaceAll('"', '\\"')}"`)
      },
    },
  ],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  server: {
    port: 5173,
    proxy: {
      "/api": "http://127.0.0.1:8080",
    },
  },
  build: {
    outDir: "../internal/handler/webui/dist",
    emptyOutDir: true,
    chunkSizeWarningLimit: 750,
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (!id.includes("node_modules")) return
          if (id.includes("@x1nx3r")) return "icons"
          if (id.includes("radix-ui")) return "radix"
          if (id.includes("@tanstack")) return "query"
          if (id.includes("recharts") || id.includes("d3-")) return "charts"
          if (/node_modules\/(react|react-dom|react-router|react-router-dom|scheduler)\//.test(id)) return "react"
          return
        },
      },
    },
  },
})
