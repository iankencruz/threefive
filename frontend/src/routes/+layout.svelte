<script lang="ts">
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { initUserContext } from '$lib/stores/user.svelte';
	import PageLoader from '$src/components/PageLoader.svelte';
	import { Toaster } from 'svelte-sonner';

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
					login({
						id: result.user.id,
						firstName: result.user.first_name,
						lastName: result.user.last_name,
						email: result.user.email,
						roles: result.user.roles ?? []
					});

					// ✅ Immediately set hydrated before redirecting
					hydrated = true;

					const pathname = page.url.pathname;
					if (pathname === '/admin/login') {
						goto('/admin/dashboard');
						return;
					}
				} else {
					logout();
					hydrated = true;

					const pathname = page.url.pathname;
					if (pathname.startsWith('/admin') && pathname !== '/admin/login') {
						goto('/admin/login');
						return;
					}
				}
			} catch (err) {
				logout();
			} finally {
				hydrated = true;
			}
		})();
	}
</script>

{#if hydrated}
	<Toaster richColors position="top-right" expand={true} />
	<main class="min-h-screen">{@render children()}</main>
{:else}
	<div class="flex h-screen w-full items-center justify-center text-gray-500">
		<PageLoader />
		<!-- <p>Loading...</p> -->
	</div>
{/if}
