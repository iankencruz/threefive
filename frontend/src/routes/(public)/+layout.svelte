<script lang="ts">
	import '$src/app.css';

	let { children, data } = $props();

	console.log('User from layout:', data.user);
</script>

<main class="flex min-h-full w-full flex-col justify-center py-12 sm:px-6 lg:px-8">
	{@render children?.()}
</main>
