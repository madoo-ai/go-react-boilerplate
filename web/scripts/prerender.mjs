import { readFile, mkdir, writeFile } from "node:fs/promises"
import { dirname, join, resolve } from "node:path"
import { fileURLToPath } from "node:url"
import { render, metaFor } from "../dist-ssr/entry-server.js"

const here = dirname(fileURLToPath(import.meta.url))
const webRoot = resolve(here, "..")
const outputDir = resolve(webRoot, "../internal/handler/webui/dist")
const indexPath = join(outputDir, "index.html")

const routes = [
  "/",
  "/login",
  "/register",
  "/forgot-password",
  "/reset-password",
  "/verify-email",
  "/two-factor",
  "/403",
]

function escapeHTML(value) {
  return value
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
}

const template = await readFile(indexPath, "utf8")

for (const route of routes) {
  const markup = render(route)
  const meta = metaFor(route)
  let html = template
    .replace('<div id="root"></div>', `<div id="root">${markup}</div>`)
    .replace(/<title>.*?<\/title>/, `<title>${escapeHTML(meta.title)}</title>`)
    .replace(
      /<meta name="description" content=".*?"\s*\/?>/,
      `<meta name="description" content="${escapeHTML(meta.description)}" />`,
    )
    .replace(/<html lang="[^"]*">/, '<html lang="id">')

  const target = route === "/" ? indexPath : join(outputDir, route.slice(1), "index.html")
  await mkdir(dirname(target), { recursive: true })
  await writeFile(target, html)
  console.log(`prerendered ${route}`)
}
