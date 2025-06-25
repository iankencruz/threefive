<script lang="ts">
	import { auth } from '$lib/store/auth.svelte';

	$effect(() => {
		if (auth.user) {
			console.log('✅ User now available:', auth.user);
		}
	});
</script>

{#if auth.user}
	<h1>Welcome back, {auth.user.first_name}!</h1>
{:else}
	<h1>Loading user...</h1>
{/if}
