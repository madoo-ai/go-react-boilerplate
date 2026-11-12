import { readFileSync } from "node:fs"

const en = JSON.parse(readFileSync(new URL("../src/messages/en.json", import.meta.url), "utf8"))
const id = JSON.parse(readFileSync(new URL("../src/messages/id.json", import.meta.url), "utf8"))

function flatten(value, prefix = "", out = new Map()) {
  if (value && typeof value === "object" && !Array.isArray(value)) {
    for (const [key, child] of Object.entries(value)) {
      flatten(child, prefix ? `${prefix}.${key}` : key, out)
    }
    return out
  }
  out.set(prefix, value)
  return out
}

const enKeys = flatten(en)
const idKeys = flatten(id)
const missingInId = [...enKeys.keys()].filter((key) => !idKeys.has(key))
const missingInEn = [...idKeys.keys()].filter((key) => !enKeys.has(key))

if (missingInId.length === 0 && missingInEn.length === 0) {
  console.log(`i18n parity OK: ${enKeys.size} keys in en and id.`)
  process.exit(0)
}

if (missingInId.length > 0) {
  console.error(`Missing in id (${missingInId.length}):`)
  for (const key of missingInId) console.error(`  ${key}`)
}
if (missingInEn.length > 0) {
  console.error(`Missing in en (${missingInEn.length}):`)
  for (const key of missingInEn) console.error(`  ${key}`)
}
process.exit(1)
