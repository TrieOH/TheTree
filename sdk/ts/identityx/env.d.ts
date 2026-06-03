/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_TRIEOH_AUTH_PROJECT_ID?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}

declare namespace NodeJS {
  interface ProcessEnv {
    readonly TRIEOH_AUTH_API_KEY?: string;
    readonly PUBLIC_TRIEOH_AUTH_PROJECT_ID?: string;
    readonly NEXT_PUBLIC_TRIEOH_AUTH_PROJECT_ID?: string;
    readonly VITE_TRIEOH_AUTH_PROJECT_ID?: string;
  }
}
