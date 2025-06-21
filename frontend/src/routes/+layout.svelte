<script>
	import { auth } from '$src/lib/store/auth.svelte.js';

	let { data, children } = $props();

	$effect(() => {
		console.log('Hydrating auth store:', data.user);
		if (data.user) {
			auth.setUser(data.user);
		} else {
			auth.clear();
		}
	});
</script>

<div class="flex min-h-full w-full flex-col justify-center py-12 sm:px-6 lg:px-8">
	{@render children?.()}
</div>
