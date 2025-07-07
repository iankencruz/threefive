<script>
	import { auth } from '$src/lib/store/auth.svelte.js';
	import { Toaster } from 'svelte-sonner';

	let { data, children } = $props();

	$effect(() => {
		console.log('Hydrating auth store:', data.user); // ✅ log here
		if (data.user) {
			auth.setUser(data.user);
		} else {
			auth.clear();
		}
	});
</script>

<Toaster closeButton richColors position="top-right" expand={true} />

<div class="flex min-h-screen w-full flex-col">
	{@render children?.()}
</div>
