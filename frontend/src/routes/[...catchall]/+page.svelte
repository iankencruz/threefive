<script lang="ts">
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';
	import { getUserContext } from '$lib/stores/user.svelte';
	import PageLoader from '$src/components/PageLoader.svelte';

	const { user } = getUserContext();
	let loaded = $state(false);

	$effect(() => {
		if (!browser) return;
		if (user.id && !loaded) {
			goto('/admin/dashboard');
			loaded = true;
		} else if (!user.id && !loaded) {
			goto('/admin/login');
			loaded = true;
		}
	});
</script>

<div class="flex h-screen flex-col items-center justify-center gap-6 text-gray-700">
	<!-- Animated Spinner -->
	<div
		class="h-12 w-12 animate-spin rounded-full border-4 border-gray-300 border-t-indigo-600"
	></div>
</div>
