import path from 'node:path';
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';

const wailsAliases = {
  '@wails': path.resolve(__dirname, 'wailsjs'),
  '@wails/go': path.resolve(__dirname, 'wailsjs/go'),
  '@wails/runtime': path.resolve(__dirname, 'wailsjs/runtime'),
};

/** Browser dev: replace Wails bindings; @wails/runtime must NOT point at wailsjs/runtime. */
const webDevAliases = {
  '@wails/go': path.resolve(__dirname, 'wailsjs/go'),
  '@wails/go/main/GUIApp': path.resolve(__dirname, 'src/dev/guiApp.ts'),
  '@wails/runtime': path.resolve(__dirname, 'src/dev'),
};

export default defineConfig(({ mode }) => {
  const isWebDev = mode === 'web';

  return {
    plugins: [react()],
    resolve: {
      alias: isWebDev ? webDevAliases : wailsAliases,
    },
    server: {
      port: 5173,
      strictPort: true,
      proxy: isWebDev
        ? {
            '/api': {
              target: 'http://localhost:8787',
              changeOrigin: true,
            },
          }
        : undefined,
    },
    build: {
      outDir: 'dist',
      emptyOutDir: true,
    },
  };
});
