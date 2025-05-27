import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig(({ mode }) => ({
  plugins: [sveltekit()],
  test: {
    include: ['src/**/*.{test,spec}.{js,ts}']
  },
  preview: {
    host: '0.0.0.0',
    port: 4173,
    allowedHosts: ['.alcatel-lucent.com', '.nokia.com', '.srexperts.net']
  },
  server: {
    proxy: mode !== 'production' ? {
      '/api': {
        target: 'http://localhost:8080/',
        rewrite: (path) => path.replace(/^\/api/, ''),
      }
    } : undefined
  }
}));
