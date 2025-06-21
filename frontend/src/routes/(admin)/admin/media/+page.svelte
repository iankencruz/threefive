<script lang="ts">
	import { auth } from '$lib/store/auth.svelte';
	console.log('Auth User first name: ', auth.user?.first_name);
</script>

{#if auth.isAuthenticated}
	<h1>Media Page: Welcome back, {auth.user?.first_name}!</h1>
	<p>isAuthenticated = {auth.isAuthenticated ? 'true' : 'false'}</p>
{:else}
	<p>You must be logged in to access this page.</p>
{/if}
