/// <reference types="vite/client" />

interface Window {
  __APP_LOCALE__?: string
  __APP_NAME__?: string
}

interface ImportMetaEnv {
  readonly VITE_APP_NAME?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
