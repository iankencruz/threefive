<script lang="ts">
	import { getUserContext } from '$lib/stores/user.svelte';
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';

	const { user } = getUserContext();
</script>

<svelte:head>
	<title>Dashboard | ThreeFive</title>
	<meta name="description" content="Login to Sabiflow" />
</svelte:head>

<h1 class="text-2xl font-bold">Welcome to your dashboard, {user.firstName}</h1>
<p class="mt-2 text-gray-600">This content is protected and only visible to logged-in users.</p>
