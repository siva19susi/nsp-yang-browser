import adapter from '@sveltejs/adapter-node';
//import adapter from '@sveltejs/adapter-cloudflare';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://kit.svelte.dev/docs/integrations#preprocessors
	// for more information about preprocessors
	preprocess: vitePreprocess(),

	kit: {
		// See https://kit.svelte.dev/docs/adapters for more information about adapters.
		adapter: adapter({
			routes: {
				include: ['/*'],
				exclude: [
					"/_app/*",
					"/fonts/*",
					"/images/*",
				]
			}
		})
	}
};

export default config;
