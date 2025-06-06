import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import path from 'path';

const config = {
	preprocess: vitePreprocess(),
	kit: {
		adapter: adapter({
			// Adjust this to wherever you want the Go backend to serve from
			pages: '../backend/static/admin',
			assets: '../backend/static/admin',
			fallback: 'index.html'
		}),
		paths: {
			base: '/admin'
		},
		alias: {
			$components: path.resolve('src/lib/components'),
			$lib: path.resolve('src/lib'),
			$stores: path.resolve('src/lib/stores'),
			$src: path.resolve('src'),
			$assets: path.resolve('src/assets')
		}
	}
};

export default config;
