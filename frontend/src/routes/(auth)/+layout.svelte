<script lang="ts">
	import '$src/app.css';
	import { goto } from '$app/navigation';
	import { initUserContext } from '$lib/stores/user.svelte';

	const { user } = initUserContext();
	let hydrated = $state(false);

	$effect(() => {
		if (hydrated && user.id !== 0) {
			goto('/admin/dashboard');
		}
	});

	if (typeof window !== 'undefined') {
		(async () => {
			hydrated = true;
		})();
	}

	let { children } = $props();
</script>

{#if hydrated}
	<main>
		{@render children()}
	</main>
{:else}
	<div class="flex h-screen items-center justify-center text-gray-500">
		Checking authentication…
	</div>
{/if}
