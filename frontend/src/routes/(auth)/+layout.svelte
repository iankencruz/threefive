<script lang="ts">
	import '$src/app.css';
	import { getUserContext } from '$lib/stores/user.svelte';
	import { goto } from '$app/navigation';
	import { browser } from '$app/environment';

	const { user, logout } = getUserContext();
	let hydrated = $state(false);

	if (browser) {
		(async () => {
			try {
				const res = await fetch('/api/v1/admin/me', { credentials: 'include' });
				const result = await res.json();

				if (res.ok && result.user?.id) {
					goto('/admin/dashboard'); // ✅ Redirect away if logged in
				} else {
					logout(); // ❗Ensure unauthenticated state is reset
					hydrated = true;
				}
			} catch (err) {
				logout(); // ❗Fail-safe
				hydrated = true;
			}
		})();
	}

	let { children } = $props();
</script>

{#if hydrated}
	<main class="flex h-screen w-full items-center justify-center text-gray-500">
		{@render children()}
	</main>
	<!-- {:else} -->
	<!-- 	<div class="flex h-screen items-center justify-center text-gray-500"> -->
	<!-- 		Checking authentication… -->
	<!-- 	</div> -->
{/if}
