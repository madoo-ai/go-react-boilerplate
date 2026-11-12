import { useEffect } from "react"
import { APP_NAME } from "@/config/app"

export function usePageMeta(title?: string, description?: string) {
  useEffect(() => {
    if (typeof document === "undefined") return
    if (title) {
      document.title = title === APP_NAME ? title : `${title} · ${APP_NAME}`
    }
    if (description) {
      let tag = document.querySelector('meta[name="description"]')
      if (!tag) {
        tag = document.createElement("meta")
        tag.setAttribute("name", "description")
        document.head.appendChild(tag)
      }
      tag.setAttribute("content", description)
    }
  }, [title, description])
}
