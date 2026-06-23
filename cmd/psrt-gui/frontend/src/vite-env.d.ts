/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_WEB_DEV?: string;
  readonly VITE_DEV_API?: string;
  readonly VITE_PRT_WEB?: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
