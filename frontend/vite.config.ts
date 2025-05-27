import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vitest/config';

export default defineConfig({
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
		proxy: {
			// svelte uses the relative /api path to reach the backend
			// can be substituted with the actual proxy like nginx
			'/api': {
				target: 'http://localhost:8080',
				rewrite: (path) => path.replace(/^\/api/, ''),
			}
		}
	}
});