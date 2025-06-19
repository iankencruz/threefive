import adapter from '@sveltejs/adapter-bun';
import { vitePreprocess } from '@sveltejs/kit/vite';

export default {
	kit: {
		adapter: adapter(),
		ssr: true
	},
	preprocess: vitePreprocess()
};
