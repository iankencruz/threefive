<script lang="ts">
	import { auth } from '$lib/store/auth.svelte';
	console.log('User', auth.user?.first_name);
</script>

<h1>Welcome back, {auth.user?.first_name}!</h1>
