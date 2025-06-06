<script lang="ts">
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';
	import { initUserContext } from '$lib/stores/user.svelte';
	import PageLoader from '$src/components/PageLoader.svelte';

	let { children } = $props();
	const { user, login, logout } = initUserContext();
	let hydrated = $state(false);

	if (browser) {
		(async () => {
			try {
				const res = await fetch('/api/v1/admin/me', {
					credentials: 'include'
				});
				const result = await res.json();

				if (res.ok && result.user?.id) {
					login(result.user);
				} else {
					logout();
				}
			} catch (err) {
				logout();
			} finally {
				setTimeout(() => {
					hydrated = true;
				}, 500); // Optional delay for smoother UX
			}
		})();
	}

	$effect(() => {
		if (!browser) return;
		if (hydrated && user.id !== 0) {
			goto('/admin/dashboard');
		}
	});
</script>

{#if hydrated}
	{console.log('✅ Hydrated and rendering children')}
	{@render children()}
{:else}
	<div class="flex h-screen items-center justify-center text-gray-500"><PageLoader /></div>
{/if}
